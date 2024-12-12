package com.ufg.SID.controller;

import com.ufg.SID.model.Lance;
import com.ufg.SID.model.Leilao;
import com.ufg.SID.model.LeilaoMensagem;
import com.ufg.SID.service.LeilaoService;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.math.BigDecimal;
import java.util.List;
import java.util.Optional;

// localhost:8080
@RestController
@RequestMapping("/auctions")
public class LeilaoController {

    @Autowired
    private RabbitTemplate rabbitTemplate;

    @Autowired
    private LeilaoService leilaoService;

    @PostMapping("/enviar")
    public String enviarMensagem(@RequestBody LeilaoMensagem mensagem) {
        rabbitTemplate.convertAndSend("leilaoQueue", mensagem);
        return "Mensagem enviada para a fila!";
    }

    // Criar leilão (POST)
    @PostMapping
    public Leilao criarLeilao(@RequestBody Leilao leilao) {
        return leilaoService.criarLeilao(leilao);
    }

    // Inscrever no leilão (POST)
    @PostMapping("/{leilaoId}/inscrever")
    public Leilao inscreverNoLeilao(@PathVariable Long leilaoId, @RequestBody Lance lance) {
        return leilaoService.inscreverNoLeilao(leilaoId, lance);
    }

    // Ver leilão específico (GET)
    @GetMapping("/{id}")
    public Optional<Leilao> verLeilao(@PathVariable Long id) {
        return leilaoService.verLeilao(id);
    }

    // Ver todos os leilões (GET)
    @GetMapping
    public List<Leilao> verTodosLeiloes() {
        return leilaoService.verTodosLeiloes();
    }

    // Ver leilões que o usuário participa (GET)
    @GetMapping("/participados")
    public List<Leilao> verLeiloesParticipados(@RequestParam String usuarioEmail) {
        return leilaoService.verLeiloesParticipados(usuarioEmail);
    }

    // Finalizar leilão (PUT)
    @PutMapping("/{id}/finalizar")
    public Leilao finalizarLeilao(@PathVariable Long id) {
        Leilao leilao = leilaoService.finalizarLeilao(id);
        LeilaoMensagem leilaoMensagem = new LeilaoMensagem();
        leilaoMensagem.setEmail(leilao.getVencedor());
        leilaoMensagem.setLeilaoProduto(leilao.getProduto());
        rabbitTemplate.convertAndSend("leilaoQueue", leilaoMensagem);
        return leilao;
    }
}

