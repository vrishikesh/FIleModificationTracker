DROP TABLE logs;

CREATE TABLE logs(path TEXT, directory TEXT, filename TEXT, inode TEXT, uid TEXT, gid TEXT, mode TEXT, device TEXT, size TEXT, block_size TEXT, atime TEXT, mtime TEXT, ctime TEXT, btime TEXT, hard_links TEXT, symlink TEXT, type TEXT);

INSERT INTO logs(path, directory, filename, inode, uid, gid, mode, device, size, block_size, atime, mtime, ctime, btime, hard_links, symlink, type) VALUES ('/Users/rishikeshvishwakarma/go/src/github.com/vrishikesh/FileModificationTracker/test.txt', '/Users/rishikeshvishwakarma/go/src/github.com/vrishikesh/FileModificationTracker', 'test.txt', '10971237', '501', '20', '0644', '0', '6', '4096', '1726066983', '1726066980', '1726066980', '1726065772', '1', '0', 'regular');
