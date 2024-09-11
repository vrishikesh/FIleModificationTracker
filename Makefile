# To install the service:
install:
	sc create FileModificationTracker binPath= "C:\path\to\binary.exe"

# To uninstall:
uninstall:
	sc delete FileModificationTracker
