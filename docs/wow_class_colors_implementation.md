# World of Warcraft Class Colors Implementation

This document details the implementation of official World of Warcraft class colors in the wclogs-cli tool using the `gookit/color` library.

## 🎨 Problem Statement

The original implementation using `fatih/color` had several issues:
1. **Limited color palette**: Only supported 16 basic ANSI colors
2. **Inaccurate colors**: Classes like Death Knight and Demon Hunter had no colors
3. **Poor distinction**: Multiple classes appeared similar (all pink/purple)
4. **Non-canonical**: Colors didn't match WoW's official class identity

## 🔍 Research & Discovery

### Official WoW Class Colors (RGB Values)

Based on WoWPedia's official RAID_CLASS_COLORS data:

| Class | RGB Values | Hex Code | Description |
|-------|------------|----------|-------------|
| Death Knight | `(196, 30, 58)` | `#C41E3A` | Deep Red |
| Demon Hunter | `(163, 48, 201)` | `#A330C9` | Purple |
| Druid | `(255, 124, 10)` | `#FF7C0A` | **Orange** |
| Hunter | `(170, 211, 114)` | `#AAD372` | Forest Green |
| Mage | `(63, 199, 235)` | `#3FC7EB` | Light Blue |
| Monk | `(0, 255, 152)` | `#00FF98` | Jade Green |
| Paladin | `(244, 140, 186)` | `#F48CBA` | **Pink** |
| Priest | `(255, 255, 255)` | `#FFFFFF` | White |
| Rogue | `(255, 244, 104)` | `#FFF468` | Yellow |
| Shaman | `(0, 112, 221)` | `#0070DD` | Blue |
| Warlock | `(135, 136, 238)` | `#8788EE` | **Light Purple** |
| Warrior | `(198, 155, 109)` | `#C69B6D` | **Brown** |
| Evoker | `(51, 147, 127)` | `#33937F` | Teal (estimated) |

### Key Discoveries
- **Druid**: Actually ORANGE `#FF7C0A`, not yellow
- **Paladin**: Pink `#F48CBA`, not white or yellow  
- **Warlock**: Light purple `#8788EE`, distinct from Demon Hunter
- **Warrior**: Brown `#C69B6D`, not red like Death Knight
- **Death Knight**: Deep red `#C41E3A`, official color

## 🛠 Technical Implementation

### Library Migration

**From:** `github.com/fatih/color` (16 ANSI colors)
**To:** `github.com/gookit/color` (RGB/TrueColor support)

```go
// Old approach (limited)
color.New(color.FgRed) // Only 16 basic colors

// New approach (RGB)
color.RGB(196, 30, 58) // Full RGB spectrum
```

### Implementation Details

```go
// getClassColor returns official WoW RGB colors
func getClassColor(class string) color.RGBColor {
    switch class {
    case "Death Knight", "DeathKnight":
        return color.RGB(196, 30, 58) // Official Deep Red
    case "Demon Hunter", "DemonHunter":
        return color.RGB(163, 48, 201) // Official Purple
    case "Druid":
        return color.RGB(255, 124, 10) // Official Orange
    case "Paladin":
        return color.RGB(244, 140, 186) // Official Pink
    case "Warlock":
        return color.RGB(135, 136, 238) // Official Light Purple
    case "Warrior":
        return color.RGB(198, 155, 109) // Official Brown
    // ... etc
    }
}
```

### Role Colors

Additionally implemented role-based colors for better visual organization:

```go
func getRoleColor(role string) color.RGBColor {
    switch role {
    case "Tank":
        return color.RGB(74, 144, 226)    // Blue
    case "Healer": 
        return color.RGB(80, 200, 120)    // Green
    case "DPS":
        return color.RGB(255, 107, 107)   // Red
    default:
        return color.RGB(153, 153, 153)   // Gray
    }
}
```

## 🎯 Results

### Before vs After

**Before (fatih/color):**
```
Warlock, Demon Hunter, Paladin → All appeared pink/purple
Death Knight → No color (white)
Druid → Yellow (incorrect)
Warrior → Same red as Death Knight
```

**After (gookit/color + official RGB):**
```
Death Knight → Deep Red #C41E3A
Demon Hunter → Purple #A330C9  
Druid → Orange #FF7C0A ✨
Hunter → Forest Green #AAD372
Mage → Light Blue #3FC7EB
Monk → Jade Green #00FF98
Paladin → Pink #F48CBA ✨
Priest → White #FFFFFF
Rogue → Yellow #FFF468
Shaman → Blue #0070DD
Warlock → Light Purple #8788EE ✨
Warrior → Brown #C69B6D ✨
```

### Terminal Output Example

```
👥 PLAYERS IN REPORT ABC123 (Fight 36) 👥
Found 20 players:

#   NAME                 CLASS        ROLE     SERVER
--- -------------------- ------------ -------- --------------------
1   Aezevera             Warlock      DPS      Antonidas
2   Akzeptabel           DeathKnight  Unknown  Eredar  
3   Athanessa            Mage         DPS      Antonidas
4   Fait                 DeathKnight  Unknown  Eredar
5   Kishun               Druid        Healer   Antonidas
6   Naalla               Paladin      Healer   Thrall
7   Sazzaria             DemonHunter  Unknown  Antonidas
8   Sharky               DemonHunter  Unknown  Sen'jin
```

Now each class is visually distinct with its authentic WoW color!

## 🔄 Migration Process

### 1. Dependency Update
```bash
go get github.com/gookit/color
```

### 2. Import Changes
```go
// Old
import "github.com/fatih/color"

// New  
import "github.com/gookit/color"
```

### 3. Function Call Updates
```go
// Old syntax
color.HiBlue("Message")
color.New(color.FgRed).Printf("Text")

// New syntax
color.Cyan.Printf("Message\n")
color.RGB(255, 0, 0).Sprint("Text")
```

### 4. Color Creation
```go
// Old (limited palette)
func getClassColor(class string) *color.Color {
    return color.New(color.FgRed)
}

// New (RGB support)  
func getClassColor(class string) color.RGBColor {
    return color.RGB(196, 30, 58)
}
```

## 🎨 Color Testing

To verify colors appear correctly in your terminal:

```bash
# Test class colors
wclogs players REPORT_CODE --debug

# Test with different terminal types
TERM=xterm-256color wclogs players REPORT_CODE
COLORTERM=truecolor wclogs players REPORT_CODE
```

## 📊 Benefits

1. **Visual Accuracy**: Matches official WoW class identity
2. **Better UX**: Each class immediately recognizable
3. **Professional**: Looks like official Blizzard tools
4. **Accessibility**: Better color distinction for users
5. **Future-proof**: RGB support allows exact color matching

## 🔮 Future Enhancements

- **Terminal Detection**: Auto-fallback for limited color terminals
- **Theme Support**: Light/dark theme variations
- **Color Customization**: User-defined color schemes
- **Accessibility**: Color-blind friendly alternatives

## 🎯 Validation

Colors validated against:
- **WoWPedia RAID_CLASS_COLORS**: Official game data
- **Battle.net**: Matches web interface colors  
- **WoW AddOns**: Consistent with popular AddOn colors
- **Terminal Testing**: Verified across multiple terminal types

The implementation now provides authentic, professional-grade World of Warcraft class colors that enhance the user experience and maintain visual consistency with the game's official color scheme.