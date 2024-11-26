package com.ufg.SID.model;

import jakarta.persistence.Embeddable;

import java.math.BigDecimal;
@Embeddable
public class Lance {

    private String usuarioEmail;
    private BigDecimal lance;

    public String getUsuarioEmail() {
        return usuarioEmail;
    }

    public void setUsuarioEmail(String usuarioEmail) {
        this.usuarioEmail = usuarioEmail;
    }

    public BigDecimal getLance() {
        return lance;
    }

    public void setLance(BigDecimal lance) {
        this.lance = lance;
    }
}
