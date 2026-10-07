# Mise cloud

Pronounced mih shih (Scottish Gaelic for "me"), Mise cloud is a developer oriented tool that allows you to easily deploy and monitor your applications on your own servers (remote or local) that are in a Tailscale tailnet. It is designed to give you an experience like Vercel, Railway, Render but on your own terms. Mise cloud does not scale to teams, and is intended for individuals (e.g. freelance or side projects)

Completion:
- CLI
    - [ ] Controller configuration
      - [ ] Controller setup
      - [ ] Auto configure tailscale services/tls
    - [ ] Deployment
    - [ ] Monitoring
- Controller
    - [ ] Agent api
    - [ ] CLI api
- Agent
  - [ ] Deployment
    - [ ] Docker management
    - [ ] Firewall management
  - [ ] Monitoring
    - [ ] Metrics collection
    - [ ] Logs streaming
    - [ ] Health checks
    - [ ] Alerts