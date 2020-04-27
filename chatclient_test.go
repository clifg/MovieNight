package main

import (
	"testing"

	"github.com/zorchenhimer/MovieNight/common"
)

func TestClient_addSpoilerTag(t *testing.T) {
	data := [][]string{
		{"||||", spoilerStart + spoilerEnd},
		{"|||||", spoilerStart + spoilerEnd + "|"},
		{"||||||", spoilerStart + spoilerEnd + "||"},
		{"|||||||", spoilerStart + spoilerEnd + "|||"},
		{"||||||||", spoilerStart + spoilerEnd + spoilerStart + spoilerEnd},
		{"||test||", spoilerStart + "test" + spoilerEnd},
		{"|| ||", spoilerStart + " " + spoilerEnd},
		{"|s|||", "|s|||"},
	}

	for i := range data {
		s := addSpoilerTags(data[i][0])
		if s != data[i][1] {
			t.Errorf("expected %#v, got %#v with %#v", data[i][1], s, data[i][0])
		}
	}
}

// Name highlighting should not interfere with emotes
func TestClient_emoteHighlight(t *testing.T) {
	client, err := NewClient(nil, nil, "Zorchenhimer", "#9547ff")
	if err != nil {
		t.Errorf("Client init error: %v", err)
	}

	data := [][]string{
		{"zorchenhimer", `<span class="mention" style="background:` + client.color + `">zorchenhimer</span>`},
		{"@zorchenhimer", `<span class="mention" style="background:` + client.color + `">@zorchenhimer</span>`},
		{"Zorchenhimer", `<span class="mention" style="background:` + client.color + `">Zorchenhimer</span>`},
		{"@Zorchenhimer", `<span class="mention" style="background:` + client.color + `">@Zorchenhimer</span>`},
		{"hello zorchenhimer", `hello <span class="mention" style="background:` + client.color + `">zorchenhimer</span>`},
		{"hello zorchenhimer ass", `hello <span class="mention" style="background:` + client.color + `">zorchenhimer</span> ass`},
		{`<img src="/emotes/twitch/zorchenhimer/zorcheWhat.png" height="40px" title="zorcheWhat">`, `<img src="/emotes/twitch/zorchenhimer/zorcheWhat.png" height="40px" title="zorcheWhat">`},
		{`zorchenhimer <img src="/emotes/twitch/zorchenhimer/zorcheWhat.png" height="40px" title="zorcheWhat">`, `<span class="mention" style="background:` + client.color + `">zorchenhimer</span> <img src="/emotes/twitch/zorchenhimer/zorcheWhat.png" height="40px" title="zorcheWhat">`},
	}

	for _, d := range data {
		chatData := client.replaceColorizedName(common.NewChatMessage(client.name, client.color, d[0], common.CmdlUser, common.MsgChat))
		if chatData.Data.(common.DataMessage).Message != d[1] {
			t.Errorf("\nExpected:\n\t%s\nReceived\n\t%s", d[1], chatData.Data.(common.DataMessage).Message)
		} else {
			t.Logf("Passed %s", d[0])
		}
	}

	// test highlighting with multiple users
	// we expect all usernames to highlight for all users
	chat = &ChatRoom{
		queue:    make(chan common.ChatData, 1),
		modqueue: make(chan common.ChatData, 1),
		clients:  []*Client{},
	}
	chat.clients = append(chat.clients, client)

	client2, err := NewClient(nil, chat, "Irani", "#9547ff")
	if err != nil {
		t.Errorf("Client init error: %v", err)
	}
	chat.clients = append(chat.clients, client2)

	data = [][]string{
		{"zorchenhimer", `<span class="othermention" style="background:` + client.color + `">zorchenhimer</span>`},
		{"@zorchenhimer", `<span class="othermention" style="background:` + client.color + `">@zorchenhimer</span>`},
		{"Zorchenhimer", `<span class="othermention" style="background:` + client.color + `">Zorchenhimer</span>`},
		{"@Zorchenhimer", `<span class="othermention" style="background:` + client.color + `">@Zorchenhimer</span>`},
		{"hello zorchenhimer", `hello <span class="othermention" style="background:` + client.color + `">zorchenhimer</span>`},
		{"hello zorchenhimer ass", `hello <span class="othermention" style="background:` + client.color + `">zorchenhimer</span> ass`},
		{"irani", `<span class="mention" style="background:` + client2.color + `">irani</span>`},
		{"@irani", `<span class="mention" style="background:` + client2.color + `">@irani</span>`},
		{`<img src="/emotes/twitch/zorchenhimer/zorcheWhat.png" height="28px" title="zorcheWhat">`, `<img src="/emotes/twitch/zorchenhimer/zorcheWhat.png" height="28px" title="zorcheWhat">`},
		{`zorchenhimer <img src="/emotes/twitch/zorchenhimer/zorcheWhat.png" height="28px" title="zorcheWhat">`, `<span class="othermention" style="background:` + client.color + `">zorchenhimer</span> <img src="/emotes/twitch/zorchenhimer/zorcheWhat.png" height="28px" title="zorcheWhat">`},
	}

	for _, d := range data {
		chatData := client2.replaceColorizedName(common.NewChatMessage(client.name, client.color, d[0], common.CmdlUser, common.MsgChat))
		if chatData.Data.(common.DataMessage).Message != d[1] {
			t.Errorf("\nExpected:\n\t%s\nReceived\n\t%s", d[1], chatData.Data.(common.DataMessage).Message)
		} else {
			t.Logf("Passed %s", d[0])
		}
	}
}
