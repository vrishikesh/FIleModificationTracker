package model

type Log struct {
	Path      string `db:"path"`
	Directory string `db:"directory"`
	Filename  string `db:"filename"`
	Inode     string `db:"inode"`
	Uid       string `db:"uid"`
	Gid       string `db:"gid"`
	Mode      string `db:"mode"`
	Device    string `db:"device"`
	Size      string `db:"size"`
	BlockSize string `db:"block_size"`
	Atime     string `db:"atime"`
	Mtime     string `db:"mtime"`
	Ctime     string `db:"ctime"`
	Btime     string `db:"btime"`
	HardLinks string `db:"hard_links"`
	Symlink   string `db:"symlink"`
	Type      string `db:"type"`
}
