# wclogs-cli

A fast, terminal-based CLI tool that wraps the Warcraft Logs GraphQL API for power users who need scriptable access to combat data without browser overhead.

## 🎯 Project Goals

**Learning Objectives:**
- Master GraphQL integration and query optimization
- Implement OAuth2 authentication flow
- Build professional CLI tools with clean UX
- Demonstrate real-world API integration skills

**Portfolio Value:**
This project showcases the ability to learn new technologies (GraphQL) and integrate with complex APIs - skills that translate directly to any company's internal systems and data platforms.

## 🚀 Features

- **Lightning Fast**: Terminal-based interface for instant data access
- **Multiple Data Types**: Damage, healing, deaths, interrupts, and more
- **Flexible Output**: ASCII tables, CSV, or JSON export
- **Smart Caching**: Respects API rate limits with intelligent caching
- **Scriptable**: Perfect for automation and data pipeline integration
- **User-Friendly**: Helpful error messages and intuitive commands

## 🏗️ Architecture

```
wclogs-cli/
├── cmd/           # CLI commands (Cobra framework)
├── api/           # GraphQL API client
├── auth/          # OAuth2 authentication
├── display/       # Terminal visualization
├── models/        # Data structures
├── config/        # Configuration handling
└── main.go        # Entry point
```

**Tech Stack:**
- **Go**: Core language for performance and CLI tooling
- **GraphQL**: Modern API integration
- **Cobra**: Professional CLI framework
- **OAuth2**: Secure authentication

## 📋 Requirements

- Go 1.19+
- Warcraft Logs API credentials (free at warcraftlogs.com)

## ⚡ Quick Start

```bash
# Install
go install github.com/yourusername/wclogs-cli@latest

# Configure (interactive setup)
wclogs config

# Get damage data for a fight
wclogs damage ABC123 5

# Export healing data to CSV
wclogs healing ABC123 5 --format csv --output healing.csv

# List all fights in a report
wclogs list ABC123
```

## 🎮 Usage Examples

```bash
# Basic damage table
wclogs damage rMGYbP9QW6KFvD4H 12

# Top 5 healers with color output
wclogs healing rMGYbP9QW6KFvD4H 12 --top 5

# Deaths for specific player
wclogs deaths rMGYbP9QW6KFvD4H 12 --player "Xaryu"

# Export interrupts as JSON
wclogs interrupts rMGYbP9QW6KFvD4H 12 --format json

# Pipe to other tools
wclogs damage ABC123 5 --format csv | head -10
```

## 🛠️ Commands

| Command | Description | Options |
|---------|-------------|---------|
| `config` | Set up API credentials interactively | |
| `damage <report> <fight>` | Show damage done table | `--top N`, `--player NAME` |
| `healing <report> <fight>` | Show healing done table | `--top N`, `--player NAME` |
| `deaths <report> <fight>` | Show death events | `--player NAME` |
| `interrupts <report> <fight>` | Show interrupt data | `--player NAME` |
| `list <report>` | List all fights in report | |

**Global Flags:**
- `--format`: Output format (`table`, `csv`, `json`)
- `--output`: Save to file instead of stdout
- `--no-cache`: Force fresh API data
- `--verbose`: Enable debug output

## ⚙️ Configuration

Create `~/.wclogs.yaml`:
```yaml
client_id: your_warcraft_logs_client_id
client_secret: your_warcraft_logs_client_secret
```

Or use the interactive setup:
```bash
wclogs config
```

## 🎯 Target Users

**Raid Leaders**: Quick damage/healing summaries between pulls without alt-tabbing
**Theorycrafters**: Export data for spreadsheet analysis and performance tracking
**Guild Officers**: Generate reports for raid performance discussions
**Developers**: Integrate Warcraft Logs data into custom tools and dashboards

## 🔧 Development Status

**Current Phase**: Week 1 - Core Functionality (Oct 18-24, 2025)
- ✅ Project scaffolding and planning
- 🔄 OAuth2 authentication implementation
- ⏳ GraphQL query engine
- ⏳ Terminal display system

**Planned Features**:
- Advanced filtering and sorting options
- Batch processing for multiple reports
- Plugin system for custom data types
- Performance analytics and insights

## 🤝 Contributing

This is primarily a learning project, but suggestions and feedback are welcome! Feel free to:
- Open issues for bugs or feature requests
- Submit PRs for improvements
- Share usage examples or tips

## 📈 Why This Project Matters

In today's data-driven development landscape, the ability to quickly integrate with GraphQL APIs and build user-friendly CLI tools is essential. This project demonstrates:

- **Rapid Technology Adoption**: Learning GraphQL from zero to production in 2 weeks
- **Real-World Problem Solving**: Addressing actual pain points for gaming communities
- **Professional Development Practices**: Clean code, documentation, testing, and user experience focus
- **API Integration Expertise**: Skills that transfer to any company's internal systems

## 📄 License

MIT License - see LICENSE file for details.

---

*Built with ❤️ as part of a structured learning journey toward professional software development.*
*"Building tomorrow's positronic brains, one CLI at a time." - Developer Learning Roadmap 2025*
