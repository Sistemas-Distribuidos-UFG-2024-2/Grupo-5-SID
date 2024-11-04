package com.ufg.SID.model;

import java.io.Serializable;

public class LeilaoMensagem implements Serializable {
    private static final long serialVersionUID = 1L;

    private Long leilaoId;
    private String email;

    public Long getLeilaoId() {
        return leilaoId;
    }

    public void setLeilaoId(Long leilaoId) {
        this.leilaoId = leilaoId;
    }

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }
}

