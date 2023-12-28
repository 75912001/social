package main

import (
	"context"
	"os"
	"path/filepath"
	blogserver "social/internal/blog"
	cleansingserver "social/internal/cleansing"
	"social/internal/friend"
	"social/internal/gate"
	interactionserver "social/internal/interaction"
	notificationserver "social/internal/notification"
	recommendationserver "social/internal/recommendation"
	"social/internal/robot"
	liblog "social/lib/log"
	pkgserver "social/pkg/server"
	"strconv"
)

func main() {
	normal := pkgserver.GetInstance()
	args := os.Args
	argNum := len(args)
	if argNum != 4 { // program name, zoneID, serviceName, serverID
		liblog.PrintfErr("args len err %v", argNum)
		return
	}
	pathName := filepath.ToSlash(args[0])
	normal.ProgramPath = filepath.ToSlash(filepath.Dir(pathName))
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
	var app pkgserver.IServer
	switch normal.ServiceName {
	case pkgserver.NameGate:
		app = gate.NewGate(normal)
	case pkgserver.NameFriend:
		app = friend.NewFriend(normal)
	case pkgserver.NameInteraction:
		app = &interactionserver.Server{Normal: normal}
	case pkgserver.NameNotification:
		app = &notificationserver.Server{Normal: normal}
	case pkgserver.NameBlog:
		app = &blogserver.Server{Normal: normal}
	case pkgserver.NameRecommendation:
		app = &recommendationserver.Server{Normal: normal}
	case pkgserver.NameCleansing:
		app = &cleansingserver.Server{Normal: normal}
	case pkgserver.NameRobot:
		app = robot.NewRobot(normal)
	default:
		liblog.PrintErr("service name err", normal.ServiceName)
		return
	}
	err := app.OnLoadBench(context.Background(), normal.Options)
	if err != nil {
		liblog.PrintErr("service name err", normal.ServiceName, err)
	}
	err = app.OnInit(context.Background(), normal.Options)
	if err != nil {
		liblog.PrintErr("service name err", normal.ServiceName, err)
		return
	}
	err = app.OnStart(context.Background())
	if err != nil {
		liblog.PrintErr("service name err", normal.ServiceName, err)
		return
	}
	err = app.OnRun(context.Background())
	if err != nil {
		liblog.PrintErr("service name err", normal.ServiceName, err)
		return
	}
	err = app.OnPreStop(context.Background())
	if err != nil {
		liblog.PrintErr("service name err", normal.ServiceName, err)
		return
	}
	err = app.OnStop(context.Background())
	if err != nil {
		liblog.PrintErr("service name err", normal.ServiceName, err)
		return
	}
	return
}
