palette = tuple((0, 0, c) for c in range(50, 256, 16)) + \
          tuple((0, c, 255) for c in range(0, 256, 6)) + \
          tuple((0, 255, 255 - c) for c in range(0, 180, 4))
for r, g, b in palette:
    print("color.RGBA{{{}, {}, {}, {}}},".format(r, g, b, 255))