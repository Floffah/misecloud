# Controller Tailscale Setup

Configure one controller service per tailnet, using the name `misecl-controller` and TCP port `50051`. This gives the controller a stable address for the CLI and agents.

## 1. Connect your devices

Install and sign in to [Tailscale](https://tailscale.com/download) on the controller machine and the machines running the CLI or agents. Use the same tailnet and enable MagicDNS. Use Tailscale 1.94 or later to avoid additional service route configuration on older Linux clients.

Use a dedicated machine for the controller. Applying its tag replaces its user identity, so user-based access rules will no longer apply to it.

## 2. Tag the controller

Open **Access controls** in the Tailscale admin console. Add this entry to your existing `tagOwners` object:

```json
"tag:misecl-controller": ["autogroup:admin"]
```

Save the policy. On the **Machines** page, find the controller device, select **Edit tags**, and apply `tag:misecl-controller`. The CLI machine and agents do not need this controller tag.

See [Tailscale's tag guide](https://tailscale.com/docs/features/tags) for tag ownership and device identity details.

## 3. Define the Tailscale Service

In the Tailscale admin console, open **Services**, select **Advertise**, then **Define a Service**. Set:

- Name: `misecl-controller` (its policy identifier is `svc:misecl-controller`).
- Endpoint: `tcp:50051`.
- Description: any label that helps you identify the controller.

The device tag and service name can both be `misecl-controller`; their `tag:` and `svc:` prefixes distinguish them.

## 4. Allow connections

Add the following rules to your existing `grants` array, preserving your other rules:

```json
{
  "src": ["autogroup:member"],
  "dst": ["tag:misecl-controller"],
  "ip": ["tcp:50051"]
},
{
  "src": ["autogroup:member"],
  "dst": ["svc:misecl-controller"],
  "ip": ["tcp:50051"]
}
```

The first rule allows direct device access for initial setup. The second allows service access. These examples allow user-owned devices belonging to tailnet members; narrow the sources to your own Tailscale login if needed. For agents on tagged devices, add their existing tags to the service rule's `src` list.

## 5. Start and configure misecloud

Start the controller on its machine. It prints its listen address and a setup command. Run that command from your CLI machine, for example:

```sh
misecl controller setup controller-host.example-tailnet.ts.net:50051
```

Replace the hostname with the controller's actual MagicDNS name. Keep the default port `50051` for this guide. If you choose another backend port, restart the controller and adjust the direct-access rule and forwarding target below.

## 6. Advertise from the controller machine

Find the controller's Tailscale IPv4 address:

```sh
tailscale ip -4
```

Run the following on that same machine, replacing `100.x.y.z` with the address returned:

```sh
tailscale serve --service=svc:misecl-controller --tcp=50051 tcp://100.x.y.z:50051
```

The current controller listens on its Tailscale IP, so forwarding to `127.0.0.1` will not reach it. The target port must match the controller's running listener. Service mode persists in the background automatically. See the [Serve command reference](https://tailscale.com/docs/reference/tailscale-cli/serve).

## 7. Approve and check the host

Return to **Services**, open `misecl-controller`, and approve the pending controller host. Host approval is separate from device tagging and access grants.

On the controller machine, check the forwarding configuration:

```sh
tailscale serve get-config --all
```

Confirm that `svc:misecl-controller` maps `tcp:50051` to the controller's actual listen address. In the console, the service should have an approved, connected host.

From another tailnet device, check that its TCP port is reachable:

```sh
nc -vz misecl-controller.example-tailnet.ts.net 50051
```

Replace `example-tailnet.ts.net` with your tailnet's MagicDNS suffix. This checks connectivity, not successful agent registration.

The standard controller address is `misecl-controller.<MagicDNS suffix>:50051`. Keep that service name for automatic agent discovery; implementing agent discovery and registration is still pending.

For more details, see [Tailscale Services](https://tailscale.com/docs/features/tailscale-services).

## Troubleshooting

- **No hosts:** confirm the controller device has `tag:misecl-controller` and advertisement ran on the controller machine.
- **Pending host:** approve it on the Services page.
- **Connection refused:** check that the controller is running and the forwarding target matches its listen address.
- **Timeout:** check host approval, access grants, and the controller's firewall.
- **Hostname does not resolve:** check MagicDNS and the client's Tailscale DNS settings.
