package hosts

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type (
	Host struct {
		ID        int    `json:"id"`
		Show      bool   `json:"show"`
		IP        string `json:"ip"`
		Hostname  string `json:"hostName"`
		GroupName string `json:"groupName"`
	}

	List map[int]*Host

	Group struct {
		HostsText string `json:"hostsText"`
		GroupName string `json:"groupName"`
		ShowNum   int    `json:"showNum"`
		HideNum   int    `json:"hideNum"`
		Show      bool   `json:"show"`
		List      List   `json:"list"`
	}

	ListByGroup map[string]Group

	MyHosts struct {
		Path             string      `json:"path"`
		HostsText        string      `json:"hostsText"`
		InUseHostsText   string      `json:"inUseHostsText"`
		NoInUseHostsText string      `json:"noInUseHostsText"`
		TotalNum         int         `json:"totalNum"`
		List             List        `json:"list"`
		ListByGroup      ListByGroup `json:"listByGroup"`
	}

	M map[string]string
)

func (m *M) toSlice() []string {
	s := make([]string, 0, len(*m))
	for _, v := range *m {
		s = append(s, v)
	}
	return s
}

func NewMyHosts() *MyHosts {
	return &MyHosts{
		Path:        getHostPath(),
		HostsText:   "",
		TotalNum:    0,
		List:        List{},
		ListByGroup: ListByGroup{},
	}
}

func (f *MyHosts) Read() error {
	iot, err := os.ReadFile(f.Path)
	if err != nil {
		return err
	}
	f.HostsText = string(iot)
	f.HostsText = strings.ReplaceAll(f.HostsText, "	", " ")
	return nil
}

func (f *MyHosts) Write() error {
	return os.WriteFile(f.Path, []byte(f.HostsText), 666)
}

func (f *MyHosts) Print() {
	log.Println(f.HostsText)
}

func (f *MyHosts) Split() {
	f.TotalNum = 0
	f.ListByGroup = ListByGroup{}
	f.List = List{}
	hosts := strings.Split(f.HostsText, "\n")
	for _, host := range hosts {
		re := regexp.MustCompile(`^([# ]*)([0-9]{0,3}\.[0-9]{0,3}\.[0-9]{0,3}\.[0-9]{0,3}|[a-zA-Z0-9:]{2,})[ ]+([a-zA-Z0-9. \-]*)+([#]+[ ]*(.*?))?$`)
		res := re.FindStringSubmatch(host)
		if len(res) == 0 {
			continue
		}
		show := !strings.Contains(res[1], "#")
		for _, hostName := range strings.Split(res[3], " ") {
			hostName = strings.TrimSpace(hostName)
			if hostName == "" {
				continue
			}
			if !strings.Contains(res[2], ":") && !strings.Contains(res[2], ".") {
				continue
			}
			var groupNames = map[string]string{}
			for _, groupName := range strings.Split(res[5], "#") {
				groupName = strings.TrimSpace(groupName)
				if groupName == "" {
					continue
				}
				groupNames[groupName] = groupName
			}
			if len(groupNames) == 0 {
				groupNames = map[string]string{"uncategorized": "uncategorized"}
			}
			for _, groupName := range groupNames {
				if _, ok := f.ListByGroup[groupName]; !ok {
					f.ListByGroup[groupName] = Group{
						HostsText: "",
						GroupName: groupName,
						ShowNum:   0,
						HideNum:   0,
						Show:      true,
						List:      map[int]*Host{},
					}
				}

				groupInfo := f.ListByGroup[groupName]
				if show {
					groupInfo.AddShowNum()
				} else {
					groupInfo.AddHideNum()
					groupInfo.Switch(false)
				}
				f.ListByGroup[groupName] = groupInfo

				f.TotalNum++
				f.ListByGroup[groupName].List[f.TotalNum] = &Host{
					ID:        f.TotalNum,
					Show:      show,
					IP:        res[2],
					Hostname:  hostName,
					GroupName: groupName,
				}
				f.List[f.TotalNum] = &Host{
					ID:        f.TotalNum,
					Show:      show,
					IP:        res[2],
					Hostname:  hostName,
					GroupName: groupName,
				}
				log.Println(Host{
					ID:        f.TotalNum,
					Show:      show,
					IP:        res[2],
					Hostname:  hostName,
					GroupName: groupName,
				})
			}
		}
	}
}

func (f *MyHosts) PrettyByGroup() {
	f.HostsText = ""
	f.InUseHostsText = ""
	f.NoInUseHostsText = ""
	for i, group := range f.ListByGroup {
		group.HostsText = ""
		for _, row := range group.List {
			if row.Show {
				group.HostsText += fmt.Sprintf("%s %s\n", row.IP, row.Hostname)
				f.HostsText += fmt.Sprintf("%s %s # %s\n", row.IP, row.Hostname, group.GroupName)
				f.InUseHostsText += fmt.Sprintf("%s %s # %s\n", row.IP, row.Hostname, group.GroupName)
			} else {
				group.HostsText += fmt.Sprintf("# %s %s\n", row.IP, row.Hostname)
				f.HostsText += fmt.Sprintf("# %s %s # %s\n", row.IP, row.Hostname, group.GroupName)
				f.NoInUseHostsText += fmt.Sprintf("# %s %s # %s\n", row.IP, row.Hostname, group.GroupName)
			}
		}
		f.ListByGroup[i] = group
	}
}

func (f *MyHosts) PrettyByList() {
	f.HostsText = ""
	f.InUseHostsText = ""
	f.NoInUseHostsText = ""
	for _, row := range f.List {
		if row.Show {
			f.HostsText += fmt.Sprintf("%s %s # %s\n", row.IP, row.Hostname, row.GroupName)
			f.InUseHostsText += fmt.Sprintf("%s %s # %s\n", row.IP, row.Hostname, row.GroupName)
		} else {
			f.HostsText += fmt.Sprintf("# %s %s # %s\n", row.IP, row.Hostname, row.GroupName)
			f.NoInUseHostsText += fmt.Sprintf("# %s %s # %s\n", row.IP, row.Hostname, row.GroupName)
		}
	}
}

func (f *MyHosts) SetGroup(groupName, hostsText string) {
	f.ListByGroup[groupName] = Group{
		HostsText: "",
		GroupName: groupName,
		ShowNum:   0,
		HideNum:   0,
		Show:      true,
		List:      map[int]*Host{},
	}
	hosts := strings.Split(hostsText, "\n")
	for _, host := range hosts {
		re := regexp.MustCompile(`^([# ]*)([0-9]{0,3}\.[0-9]{0,3}\.[0-9]{0,3}\.[0-9]{0,3}|[a-zA-Z0-9:]{2,})[ ]+([a-zA-Z0-9. \-]*)+([#]+[ ]*(.*?))?$`)
		res := re.FindStringSubmatch(host)
		if len(res) == 0 {
			continue
		}
		show := !strings.Contains(res[1], "#")
		for _, hostName := range strings.Split(res[3], " ") {
			hostName = strings.TrimSpace(hostName)
			if hostName == "" {
				continue
			}
			if !strings.Contains(res[2], ":") && !strings.Contains(res[2], ".") {
				continue
			}
			f.TotalNum++
			f.ListByGroup[groupName].List[f.TotalNum] = &Host{
				ID:        f.TotalNum,
				Show:      show,
				IP:        res[2],
				Hostname:  hostName,
				GroupName: groupName,
			}
		}
	}

	groupInfo := f.ListByGroup[groupName]
	groupInfo.Switch(true)
	for _, row := range groupInfo.List {
		if row.Show {
			groupInfo.AddShowNum()
		} else {
			groupInfo.Switch(false)
			groupInfo.AddHideNum()
		}
	}
	f.ListByGroup[groupName] = groupInfo
}

func (f *MyHosts) Add(groupName string, ip string, hostName string) {
	f.TotalNum++
	if _, ok := f.ListByGroup[groupName]; !ok {
		f.ListByGroup[groupName] = Group{
			HostsText: "",
			GroupName: groupName,
			ShowNum:   0,
			HideNum:   0,
			Show:      true,
			List:      map[int]*Host{},
		}
	}
	f.ListByGroup[groupName].List[f.TotalNum] = &Host{
		ID:        f.TotalNum,
		Show:      true,
		IP:        ip,
		Hostname:  hostName,
		GroupName: groupName,
	}
}

func (f *MyHosts) Delete(groupName string, hostNameID int) {
	if _, ok := f.ListByGroup[groupName]; !ok {
		return
	}
	delete(f.ListByGroup[groupName].List, hostNameID)
	if len(f.ListByGroup[groupName].List) == 0 {
		delete(f.ListByGroup, groupName)
	}
}

func (f *MyHosts) SwitchByGroupName(groupName string, show bool) {
	if _, ok := f.ListByGroup[groupName]; !ok {
		return
	}
	groupInfo := f.ListByGroup[groupName]
	groupInfo.Switch(show)
	for id := range groupInfo.List {
		groupInfo.List[id].Switch(show)
		f.List[id].Switch(show)
	}
	if show {
		groupInfo.ShowNum = len(groupInfo.List)
		groupInfo.HideNum = 0
	} else {
		groupInfo.ShowNum = 0
		groupInfo.HideNum = len(groupInfo.List)
	}
	f.ListByGroup[groupName] = groupInfo
}

func (f *MyHosts) SwitchByHostNameId(groupName string, hostNameID int, show bool) {
	if _, ok := f.ListByGroup[groupName]; !ok {
		return
	}
	groupInfo := f.ListByGroup[groupName]
	if _, ok := groupInfo.List[hostNameID]; !ok {
		return
	}
	groupInfo.List[hostNameID].Switch(show)
	groupInfo.Switch(true)
	groupInfo.ShowNum = 0
	groupInfo.HideNum = 0
	for _, row := range groupInfo.List {
		if row.Show {
			groupInfo.AddShowNum()
		} else {
			groupInfo.Switch(false)
			groupInfo.AddHideNum()
		}
	}
	f.ListByGroup[groupName] = groupInfo
	f.List[hostNameID].Switch(show)
}

func (f *MyHosts) SetGroupNameByOldGroupName(oldGroupName, groupName string) {
	if _, ok := f.ListByGroup[oldGroupName]; !ok {
		return
	}
	oldGroupInfo := f.ListByGroup[oldGroupName]
	delete(f.ListByGroup, oldGroupName)

	if _, ok := f.ListByGroup[groupName]; !ok {
		f.ListByGroup[groupName] = Group{
			HostsText: "",
			GroupName: groupName,
			ShowNum:   0,
			HideNum:   0,
			Show:      true,
			List:      map[int]*Host{},
		}
	}

	for _, row := range oldGroupInfo.List {
		f.ListByGroup[groupName].List[row.ID] = &Host{
			ID:        row.ID,
			Show:      row.Show,
			IP:        row.IP,
			Hostname:  row.Hostname,
			GroupName: groupName,
		}
	}

	groupInfo := f.ListByGroup[groupName]
	groupInfo.Switch(true)
	groupInfo.ShowNum = 0
	groupInfo.HideNum = 0
	for _, row := range groupInfo.List {
		if row.Show {
			groupInfo.AddShowNum()
		} else {
			groupInfo.Switch(false)
			groupInfo.AddHideNum()
		}
	}
	f.ListByGroup[groupName] = groupInfo
}

func (f *MyHosts) GetAllGroupNames() []string {
	var groupNames []string
	for _, g := range f.ListByGroup {
		groupNames = append(groupNames, g.GroupName)
	}
	return groupNames
}

func (f *MyHosts) SetGroupNameByHostnameId(hostNameID int, groupName string) {
	if _, ok := f.List[hostNameID]; !ok {
		return
	}
	f.List[hostNameID].SetGroupName(groupName)
}

func (h *Host) Switch(show bool) {
	h.Show = show
}

func (h *Host) SetGroupName(groupName string) {
	h.GroupName = groupName
}

func (g *Group) Switch(show bool) {
	g.Show = show
}

func (g *Group) AddHideNum() {
	g.HideNum++
}

func (g *Group) AddShowNum() {
	g.ShowNum++
}
