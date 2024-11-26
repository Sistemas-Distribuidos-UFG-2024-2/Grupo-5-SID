package com.ufg.SID.service;

import com.ufg.SID.model.Lance;
import com.ufg.SID.model.Leilao;
import com.ufg.SID.repository.LeilaoRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;


import java.util.Comparator;
import java.util.List;
import java.util.Optional;

@Service
public class LeilaoService {

    @Autowired
    private LeilaoRepository leilaoRepository;

    public Leilao criarLeilao(Leilao leilao) {
        return leilaoRepository.save(leilao);
    }

    public Optional<Leilao> verLeilao(Long id) {
        return leilaoRepository.findById(id);
    }

    public List<Leilao> verTodosLeiloes() {
        return leilaoRepository.findAll();
    }

    public List<Leilao> verLeiloesParticipados(String usuarioEmail) {
        return leilaoRepository.findByUsuarioParticipante(usuarioEmail);
    }

    public Leilao inscreverNoLeilao(Long leilaoId, Lance lance) {
        Leilao leilao = leilaoRepository.findById(leilaoId).orElseThrow(() -> new RuntimeException("Leilão não encontrado"));
        leilao.getParticipantes().add(lance);
        return leilaoRepository.save(leilao);
    }

    public Leilao finalizarLeilao(Long leilaoId) {
        // Procura o leilão pelo ID
        Leilao leilao = leilaoRepository.findById(leilaoId)
                .orElseThrow(() -> new RuntimeException("Leilão não encontrado"));

        if (leilao.isFinalizado()) {
            throw new RuntimeException("Leilão já finalizado.");
        }

        // Participante com o maior lance
        Lance vencedor = leilao.getParticipantes().stream()
                .max(Comparator.comparing(Lance::getLance))
                .orElseThrow(() -> new RuntimeException("Nenhum lance foi registrado."));

        leilao.setVencedor(vencedor.getUsuarioEmail());
        leilao.setLanceFinal(vencedor.getLance());
        leilao.setFinalizado(true);

        return leilaoRepository.save(leilao);
    }
}
