package main

import (
	"context"
	"os"
	"path/filepath"
	blogserver "social/internal/blog"
	cleansingserver "social/internal/cleansing"
	friendserver "social/internal/friend"
	gateserver "social/internal/gate/server"
	interactionserver "social/internal/interaction"
	notificationserver "social/internal/notification"
	recommendationserver "social/internal/recommendation"
	robotserver "social/internal/robot/server"
	"social/lib/log"
	libruntime "social/lib/runtime"

	liblog "social/lib/log"
	//pkgcommon "social/pkg/common"
	pkgserver "social/pkg/server"
	"strconv"
)

func main() {
	normal := pkgserver.GetInstance()
	args := os.Args
	argNum := len(args)
	if argNum != 4 { // program name, zoneID, serviceName, serverID
		log.PrintErr("args len err")
		return
	}
	pathName := filepath.ToSlash(args[0])
	normal.ProgramPath = filepath.Dir(pathName)
	normal.ProgramName = filepath.Base(pathName)
	{
		strZoneID, err := strconv.ParseUint(args[1], 10, 32)
		if err != nil {
			log.PrintErr("zoneID err", err)
			return
		}
		normal.ZoneID = uint32(strZoneID)
	}
	normal.ServiceName = args[2]
	{
		strServiceID, err := strconv.ParseUint(args[3], 10, 32)
		if err != nil {
			log.PrintErr("serviceID err", err)
			return
		}
		normal.ServiceID = uint32(strServiceID)
	}
	log.PrintInfo(normal.ZoneID, normal.ServiceName, normal.ServiceID)
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
		s = robotserver.NewServer(normal)
	default:
		log.PrintErr("service name err", normal.ServiceName)
		return
	}
	err := s.OnLoadBench(context.Background(), normal.Options)
	if err != nil {
		log.PrintErr("service name err", normal.ServiceName, err)
	}
	err = s.OnInit(context.Background(), normal.Options)
	if err != nil {
		log.PrintErr("service name err", normal.ServiceName, err)
		return
	}
	liblog.GetInstance().Tracef("xxxxxxxxxxxxxxxx", libruntime.Location())
	err = s.OnStart(context.Background())
	if err != nil {
		log.PrintErr("service name err", normal.ServiceName, err)
		return
	}
	err = s.OnRun(context.Background())
	if err != nil {
		log.PrintErr("service name err", normal.ServiceName, err)
		return
	}
	err = s.OnPreStop(context.Background())
	if err != nil {
		log.PrintErr("service name err", normal.ServiceName, err)
		return
	}
	err = s.OnStop(context.Background())
	if err != nil {
		log.PrintErr("service name err", normal.ServiceName, err)
		return
	}
	return
}
