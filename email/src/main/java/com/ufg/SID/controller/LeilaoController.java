package com.ufg.SID.controller;

import com.ufg.SID.model.LeilaoMensagem;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/leilao")
public class LeilaoController {

    @Autowired
    private RabbitTemplate rabbitTemplate;

    @PostMapping("/enviar")
    public String enviarMensagem(@RequestBody LeilaoMensagem mensagem) {
        rabbitTemplate.convertAndSend("leilaoQueue", mensagem);
        return "Mensagem enviada para a fila!";
    }
}

