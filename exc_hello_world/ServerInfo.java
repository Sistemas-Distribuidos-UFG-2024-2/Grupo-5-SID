package hello_world;

public class ServerInfo {
    private String ip;
    private int port;
    private boolean enabled;

    public ServerInfo() {}

    // Construtor
    public ServerInfo(String ip, int port, boolean enabled) {
        this.ip = ip;
        this.port = port;
        this.enabled = enabled;
    }

    // Getters e Setters
    public String getIp() {
        return ip;
    }

    public void setIp(String ip) {
        this.ip = ip;
    }

    public int getPort() {
        return port;
    }

    public void setPort(int port) {
        this.port = port;
    }

    public boolean isEnabled() {
        return enabled;
    }

    public void setEnabled(boolean enabled) {
        this.enabled = enabled;
    }

    // Método toJSON
    public String toJSON() {
        return String.format(
                "{\"ip\": \"%s\", \"port\": %d, \"enabled\": %b}",
                this.ip, this.port, this.enabled
        );
    }

    // Método para exibir informações do servidor
    @Override
    public String toString() {
        return "Server { " +
                "IP='" + ip + '\'' +
                ", Port='" + port + '\'' +
                ", Enabled=" + enabled +
                " }";
    }
}
