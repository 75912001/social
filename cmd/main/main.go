package main

import (
	"context"
	"os"
	"social/internal/blog"
	"social/internal/cleansing"
	"social/internal/friend"
	"social/internal/gate"
	"social/internal/interaction"
	"social/internal/notification"
	"social/internal/recommendation"
	"social/internal/robot"
	"social/pkg"
	xrlog "social/pkg/lib/log"
	"social/pkg/server"
	"strconv"
)

func main() {
	args := os.Args
	argNum := len(args)
	if argNum != 4 { // program name, zoneID, serviceName, serverID
		xrlog.PrintErr("args len err")
		return
	}
	{
		zoneID, err := strconv.ParseUint(args[1], 10, 32)
		if err != nil {
			xrlog.PrintErr("zoneID err", err)
			return
		}
		pkg.GZoneID = uint32(zoneID)
	}
	pkg.GServiceName = args[2]
	{
		serviceID, err := strconv.ParseUint(args[3], 10, 32)
		if err != nil {
			xrlog.PrintErr("serviceID err", err)
			return
		}
		pkg.GServiceID = uint32(serviceID)
	}
	xrlog.PrintInfo(pkg.GZoneID, pkg.GServiceName, pkg.GServiceID)
	switch pkg.GServiceName {
	case server.NameGate:
		server.GetInstance().Server = &gate.Server{}
	case server.NameFriend:
		server.GetInstance().Server = &friend.Server{}
	case server.NameInteraction:
		server.GetInstance().Server = &interaction.Server{}
	case server.NameNotification:
		server.GetInstance().Server = &notification.Server{}
	case server.NameBlog:
		server.GetInstance().Server = &blog.Server{}
	case server.NameRecommendation:
		server.GetInstance().Server = &recommendation.Server{}
	case server.NameCleansing:
		server.GetInstance().Server = &cleansing.Server{}
	case server.NameRobot:
		server.GetInstance().Server = &robot.Server{}
	default:
		xrlog.PrintErr("service name err", pkg.GServiceName)
		return
	}
	err := server.GetInstance().Server.Start(context.Background())
	if err != nil {
		xrlog.PrintErr("service name err", pkg.GServiceName, err)
	}
	err = server.GetInstance().Server.Stop(context.Background())
	if err != nil {
		xrlog.PrintErr("service name err", pkg.GServiceName, err)
	}
	return
}
