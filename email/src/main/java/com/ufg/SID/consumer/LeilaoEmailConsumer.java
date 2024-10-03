package com.ufg.SID.consumer;

import com.ufg.SID.model.LeilaoMensagem;
import com.ufg.SID.service.EmailService;
import com.rabbitmq.client.Channel;
import org.springframework.amqp.core.Message;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

@Component
public class LeilaoEmailConsumer {

    @Autowired
    private EmailService emailService;

    @RabbitListener(queues = "leilaoQueue")
    public void receberMensagem(LeilaoMensagem mensagem, Channel channel, Message message) {
        try {
            emailService.enviarEmailGanhador(mensagem.getEmail(), mensagem.getLeilaoId());

            // Checa se rabbitmq retorna ack (MENSAGEM ENTREGUE COM SUCESSO)
            channel.basicAck(message.getMessageProperties().getDeliveryTag(), false);
        } catch (Exception e) {
            e.printStackTrace();
            try {
                // Caso não, reencaminha pra fila do rabbitmq pra tentar novamente
                channel.basicNack(message.getMessageProperties().getDeliveryTag(), false, true);
            } catch (Exception ex) {
                ex.printStackTrace();
            }
        }
    }
}


