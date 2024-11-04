package com.ufg.SID.service;

import com.sendgrid.*;
import com.sendgrid.helpers.mail.Mail;
import com.sendgrid.helpers.mail.objects.Content;
import com.sendgrid.helpers.mail.objects.Email;
import org.springframework.stereotype.Service;

import java.io.IOException;

@Service
public class EmailService {

    private final String apiKey = "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";

    public void enviarEmailGanhador(String email, Long leilaoId) throws IOException {
        Email from = new Email("leilaosid@gmail.com");
        String subject = "Parabéns! Você foi o ganhador!";
        Email to = new Email(email);
        Content content = new Content("text/plain", "Parabéns! Você foi o ganhador do leilão com id: " + leilaoId);
        Mail mail = new Mail(from, subject, to, content);

        SendGrid sg = new SendGrid(apiKey);
        Request request = new Request();
        try {
            request.setMethod(Method.POST);
            request.setEndpoint("mail/send");
            request.setBody(mail.build());
            Response response = sg.api(request);
            System.out.println(response.getStatusCode());
            System.out.println(response.getBody());
            System.out.println(response.getHeaders());
        } catch (IOException ex) {
            throw ex;
        }
    }
}
