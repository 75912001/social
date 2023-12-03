package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	blogserver "social/internal/blog"
	cleansingserver "social/internal/cleansing"
	friendserver "social/internal/friend"
	gateserver "social/internal/gate/server"
	interactionserver "social/internal/interaction"
	notificationserver "social/internal/notification"
	recommendationserver "social/internal/recommendation"
	robotserver "social/internal/robot"
	pkgcommon "social/pkg/common"
	liblog "social/pkg/lib/log"
	pkgserver "social/pkg/server"
	"strconv"
)

func main() {
	normal := pkgserver.NewNormal()
	args := os.Args
	argNum := len(args)
	if argNum != 4 { // program name, zoneID, serviceName, serverID
		liblog.PrintErr("args len err")
		return
	}
	pathName := filepath.ToSlash(args[0])
	normal.ProgramPath = filepath.Dir(pathName)
	normal.ProgramName = filepath.Base(pathName)
	{
		strZoneID, err := strconv.ParseUint(args[1], 10, 32)
		if err != nil {
			liblog.PrintErr("zoneID err", err)
			return
		}
		normal.ZoneID = uint32(strZoneID)
	}
	normal.ServiceName = args[2]
	{
		strServiceID, err := strconv.ParseUint(args[3], 10, 32)
		if err != nil {
			liblog.PrintErr("serviceID err", err)
			return
		}
		normal.ServiceID = uint32(strServiceID)
	}
	liblog.PrintInfo(normal.ZoneID, normal.ServiceName, normal.ServiceID)
	var s pkgserver.IServer
	switch normal.ServiceName {
	case pkgserver.NameGate:
		s = gateserver.NewServer(normal)
	case pkgserver.NameFriend:
		s = &friendserver.Server{Normal: normal}
	case pkgserver.NameInteraction:
		s = &interactionserver.Server{Normal: normal}
	case pkgserver.NameNotification:
		s = &notificationserver.Server{Normal: normal}
	case pkgserver.NameBlog:
		s = &blogserver.Server{Normal: normal}
	case pkgserver.NameRecommendation:
		s = &recommendationserver.Server{Normal: normal}
	case pkgserver.NameCleansing:
		s = &cleansingserver.Server{Normal: normal}
	case pkgserver.NameRobot:
		s = &robotserver.Server{Normal: normal}
	default:
		liblog.PrintErr("service name err", normal.ServiceName)
		return
	}
	err := s.LoadBench(context.Background(), normal.Options)
	if err != nil {
		liblog.PrintErr("service name err", normal.ServiceName, err)
	}
	err = s.Init(context.Background(),
		normal.Options, pkgserver.NewOptions().SetEtcdWatchServicePrefix(fmt.Sprintf("/%v/%v/", pkgcommon.ProjectName, pkgcommon.EtcdWatchMsgTypeService)).
			SetEtcdWatchCommandPrefix(fmt.Sprintf("/%v/%v/%v/%v/",
				pkgcommon.ProjectName, pkgcommon.EtcdWatchMsgTypeCommand,
				normal.ZoneID,
				normal.ServiceName),
			),
	)
	if err != nil {
		liblog.PrintErr("service name err", normal.ServiceName, err)
		return
	}
	err = s.Start(context.Background())
	if err != nil {
		liblog.PrintErr("service name err", normal.ServiceName, err)
		return
	}
	err = s.Run(context.Background())
	if err != nil {
		liblog.PrintErr("service name err", normal.ServiceName, err)
		return
	}
	err = s.PreStop(context.Background())
	if err != nil {
		liblog.PrintErr("service name err", normal.ServiceName, err)
		return
	}
	err = s.Stop(context.Background())
	if err != nil {
		liblog.PrintErr("service name err", normal.ServiceName, err)
		return
	}
	return
}
