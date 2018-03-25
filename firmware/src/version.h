#ifndef VERSION_H__
#define VERSION_H__

#define VER_MAJOR(x) ((x & 0x1f) << 11) // [15:11]
#define VER_MINOR(x) ((x & 0x1f) << 6)  // [10:6]
#define VER_PATCH(x) (x & 0x2f)         // [5:0]

#define FW_VERSION_(ma, mi, p) (VER_MAJOR(ma) | VER_MINOR(mi) | VER_PATCH(p))

#define FW_VERSION FW_VERSION_(0, 1, 0)

#endif

