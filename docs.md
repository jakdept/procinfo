# Documentation for Kernel interfaces used

From `man 5 proc`:
```man
/proc/locks
  This file shows current file locks (flock(2) and fcntl(2)) and leases (fcntl(2)).

  An example of the content shown in this file is the following:
      0  1      2         3     4    5             6   7
      1: POSIX  ADVISORY  READ  5433 08:01:7864448 128 128
      2: FLOCK  ADVISORY  WRITE 2001 08:01:7864554 0 EOF
      3: FLOCK  ADVISORY  WRITE 1568 00:2f:32388 0 EOF
      4: POSIX  ADVISORY  WRITE 699 00:16:28457 0 EOF
      5: POSIX  ADVISORY  WRITE 764 00:16:21448 0 0
      6: POSIX  ADVISORY  READ  3548 08:01:7867240 1 1
      7: POSIX  ADVISORY  READ  3548 08:01:7865567 1826 2335
      8: OFDLCK ADVISORY  WRITE -1 08:01:8713209 128 191

  The fields shown in each line are as follows:

  (1) The ordinal position of the lock in the list.

  (2) The lock type.  Values that may appear here include:

      FLOCK  This is a BSD file lock created using flock(2).

      OFDLCK This is an open file description (OFD) lock created using fcntl(2).

      POSIX  This is a POSIX byte-range lock created using fcntl(2).

  (3) Among the strings that can appear here are the following:

      ADVISORY
              This is an advisory lock.

      MANDATORY
              This is a mandatory lock.

  (4) The type of lock.  Values that can appear here are:

      READ   This is a POSIX or OFD read lock, or a BSD shared lock.

      WRITE  This is a POSIX or OFD write lock, or a BSD exclusive lock.

  (5) The PID of the process that owns the lock.

      Because OFD locks are not owned by a single process (since multiple processes may have file descriptors  that  refer to  the  same open file description), the value -1 is displayed in this field for OFD locks.  (Before kernel 4.14, a bug meant that the PID of the process that initially acquired the lock was displayed instead of the value -1.)

  (6) Three colon-separated subfields that identify the major and minor device ID of the device containing the  filesystem
      where the locked file resides, followed by the inode number of the locked file.

  (7) The byte offset of the first byte of the lock.  For BSD locks, this value is always 0.

  (8) The byte offset of the last byte of the lock.  EOF in this field means that the lock extends to the end of the file.
      For BSD locks, the value shown is always EOF.

  Since Linux 4.9, the list of locks shown in /proc/locks is filtered to show just the locks for the processes in the  PID
  namespace  (see  pid_namespaces(7)) for which the /proc filesystem was mounted.  (In the initial PID namespace, there is
  no filtering of the records shown in this file.)

  The lslocks(8) command provides a bit more information about each lock.
```