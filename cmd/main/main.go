package main

import (
	"os"
	"social/cmd/blog"
	"social/cmd/cleansing"
	"social/cmd/friend"
	"social/cmd/gate"
	"social/cmd/interaction"
	"social/cmd/notification"
	"social/cmd/recommendation"
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
	var s server.IServer
	switch pkg.GServiceName {
	case server.NameGate:
		s = &gate.Server{}
	case server.NameFriend:
		s = &friend.Server{}
	case server.NameInteraction:
		s = &interaction.Server{}
	case server.NameNotification:
		s = &notification.Server{}
	case server.NameBlog:
		s = &blog.Server{}
	case server.NameRecommendation:
		s = &recommendation.Server{}
	case server.NameCleansing:
		s = &cleansing.Server{}
	default:
		xrlog.PrintErr("service name err", pkg.GServiceName)
		return
	}
	err := s.Start()
	if err != nil {
		xrlog.PrintErr("service name err", pkg.GServiceName, err)
	}
	err = s.Stop()
	if err != nil {
		xrlog.PrintErr("service name err", pkg.GServiceName, err)
	}
	return
}
