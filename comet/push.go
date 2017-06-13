package main

type Push int

func (p *Push) Broad(args *Msg, reply *Reply) error {
	room := args.RoomId
	data := preMsg(args.Username, args.Content)
	for _, v := range cmap {
		if v.roomId == room {
			v.wch <- data
		}
	}
	return nil
}
