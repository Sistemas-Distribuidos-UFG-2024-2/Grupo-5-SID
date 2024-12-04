package com.ufg.SID.model;

import jakarta.persistence.*;
import java.util.UUID;

@Entity
@Table(name = "accounts", schema = "schema_accounts") // Tabela e schema correspondentes
public class Usuarios {

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO) // Geração automática para UUID
    private UUID id; // Alinhado ao campo 'id' do banco

    @Column(name = "mail", nullable = false, unique = true) // Alinhado ao campo 'mail' do banco
    private String mail;

    @Column(name = "enabled", nullable = false) // Representa o campo 'enabled' do banco
    private Boolean enabled;

    // Getters e Setters
    public UUID getId() {
        return id;
    }

    public void setId(UUID id) {
        this.id = id;
    }

    public String getMail() {
        return mail;
    }

    public void setMail(String mail) {
        this.mail = mail;
    }

    public Boolean getEnabled() {
        return enabled;
    }

    public void setEnabled(Boolean enabled) {
        this.enabled = enabled;
    }
}
