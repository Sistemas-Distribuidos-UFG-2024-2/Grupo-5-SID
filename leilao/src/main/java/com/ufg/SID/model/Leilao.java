package com.ufg.SID.model;

import jakarta.persistence.*;
import java.math.BigDecimal;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

@Entity
public class Leilao {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private String produto;
    private BigDecimal lanceInicial;
    private LocalDateTime dataFinalizacao;
    private boolean finalizado = false;
    private String criador;
    private String vencedor;
    private BigDecimal lanceFinal;
    private BigDecimal valorMaximo;

    @ElementCollection
    @CollectionTable(
            name = "participantes",
            joinColumns = @JoinColumn(name = "leilao_id")
    )
    private List<Lance> participantes = new ArrayList<>();

    public String getCriador() {
        return criador;
    }

    public void setCriador(String criador) {
        this.criador = criador;
    }

    public String getVencedor() {
        return vencedor;
    }

    public void setVencedor(String vencedor) {
        this.vencedor = vencedor;
    }

    public BigDecimal getLanceFinal() {
        return lanceFinal;
    }

    public void setLanceFinal(BigDecimal lanceFinal) {
        this.lanceFinal = lanceFinal;
    }

    public BigDecimal getValorMaximo() {
        return valorMaximo;
    }

    public void setValorMaximo(BigDecimal valorMaximo) {
        this.valorMaximo = valorMaximo;
    }

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

    public LocalDateTime getDataFinalizacao() {
        return dataFinalizacao;
    }

    public void setDataFinalizacao(LocalDateTime dataFinalizacao) {
        this.dataFinalizacao = dataFinalizacao;
    }

    public boolean isFinalizado() {
        return finalizado;
    }

    public void setFinalizado(boolean finalizado) {
        this.finalizado = finalizado;
    }

    public List<Lance> getParticipantes() {
        return participantes;
    }

    public void setParticipantes(List<Lance> participantes) {
        this.participantes = participantes;
    }
}
