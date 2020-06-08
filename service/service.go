/**
 *     ______                 __
 *    /\__  _\               /\ \
 *    \/_/\ \/     ___     __\ \ \         __      ___     ___     __
 *       \ \ \    / ___\ / __ \ \ \  __  / __ \  /  _  \  / ___\ / __ \
 *        \_\ \__/\ \__//\  __/\ \ \_\ \/\ \_\ \_/\ \/\ \/\ \__//\  __/
 *        /\_____\ \____\ \____\\ \____/\ \__/ \_\ \_\ \_\ \____\ \____\
 *        \/_____/\/____/\/____/ \/___/  \/__/\/_/\/_/\/_/\/____/\/____/
 *
 *
 *                                                                    @寒冰
 *                                                            www.icezzz.cn
 *                                                     hanbin020706@163.com
 */
package service

type Service interface {
	Start()
	Stop()
	GetServiceName() string
}

type Entity struct {
	Name string `toml:"name"`
	Ip   string `toml:"ip"`
	Port int    `toml:"port"`
}

func (s *Entity) Start() {
	panic("implement me")
}

func (s *Entity) Stop() {
	panic("implement me")
}

func (s *Entity) GetServiceName() string {
	return s.Name
}
