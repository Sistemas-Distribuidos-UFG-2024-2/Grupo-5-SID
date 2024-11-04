package com.ufg.SID.model;

import jakarta.persistence.*;
import java.math.BigDecimal;
import java.time.LocalDate;
import java.util.HashSet;
import java.util.Set;

@Entity
public class Leilao {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private String produto;
    private BigDecimal lanceInicial;
    private LocalDate dataFinalizacao;
    private boolean finalizado = false;

    // Conjunto de IDs dos usuários participantes do leilão
    @ElementCollection
    private Set<Long> participantes = new HashSet<>();

    // Getters e Setters
    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getProduto() {
        return produto;
    }

    public void setProduto(String produto) {
        this.produto = produto;
    }

    public BigDecimal getLanceInicial() {
        return lanceInicial;
    }

    public void setLanceInicial(BigDecimal lanceInicial) {
        this.lanceInicial = lanceInicial;
    }

    public LocalDate getDataFinalizacao() {
        return dataFinalizacao;
    }

    public void setDataFinalizacao(LocalDate dataFinalizacao) {
        this.dataFinalizacao = dataFinalizacao;
    }

    public boolean isFinalizado() {
        return finalizado;
    }

    public void setFinalizado(boolean finalizado) {
        this.finalizado = finalizado;
    }

    public Set<Long> getParticipantes() {
        return participantes;
    }

    public void setParticipantes(Set<Long> participantes) {
        this.participantes = participantes;
    }
}
