package com.ufg.SID.service;

import com.ufg.SID.model.Leilao;
import com.ufg.SID.repository.LeilaoRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

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

    public List<Leilao> verLeiloesParticipados(Long usuarioId) {
        return leilaoRepository.findByParticipantesContains(usuarioId);
    }

    public Leilao inscreverNoLeilao(Long leilaoId, Long usuarioId) {
        Leilao leilao = leilaoRepository.findById(leilaoId).orElseThrow(() -> new RuntimeException("Leilão não encontrado"));
        leilao.getParticipantes().add(usuarioId);
        return leilaoRepository.save(leilao);
    }

    public Leilao finalizarLeilao(Long leilaoId) {
        Leilao leilao = leilaoRepository.findById(leilaoId).orElseThrow(() -> new RuntimeException("Leilão não encontrado"));
        leilao.setFinalizado(true);
        return leilaoRepository.save(leilao);
    }
}
