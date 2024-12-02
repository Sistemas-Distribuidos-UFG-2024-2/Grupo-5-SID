package com.ufg.SID.model;

import java.io.Serializable;

public class LeilaoMensagem implements Serializable {
    private static final long serialVersionUID = 1L;

    private String leilaoProduto;
    private String email;

    public String getLeilaoProduto() {
        return leilaoProduto;
    }

    public void setLeilaoProduto(String leilaoProduto) {
        this.leilaoProduto = leilaoProduto;
    }

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }
}

