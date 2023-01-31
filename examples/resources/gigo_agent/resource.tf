data "gigo_workspace" "me" {
}

resource "gigo_agent" "dev" {
  os   = "linux"
  arch = "amd64"
  dir  = "/workspace"
}

resource "kubernetes_pod" "dev" {
  count = data.gigo_workspace.me.start_count
  spec {
    container {
      command = ["sh", "-c", gigo_agent.dev.init_script]
      env {
        name  = "GIGO_AGENT_TOKEN"
        value = gigo_agent.dev.token
      }
    }
  }
}
