import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '30s', target: 1000 },  // Aumenta de 0 para 10 usuários simultâneos em 30 segundos
        { duration: '1m', target: 2000 },   // Aumenta para 50 usuários em 1 minuto
        { duration: '30s', target: 0 },   // Diminui de 50 para 0 usuários em 30 segundos
    ],
};

export default function () {
    // URL da rota
    const url = 'http://localhost:32300/auctions/2/bids';

    // Payload com dados para o POST
    const payload = JSON.stringify({
        account_id: "1",
        amount: Math.floor(Math.random() * 100) + 1,  // Aumenta o valor de "amount" (pode ser ajustado conforme necessidade)
    });

    // Configurações da requisição
    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const response = http.post(url, payload, params);

    check(response, {
        'status is 202': (r) => r.status === 202,
    });

    // Atraso entre requisições (pode ser ajustado conforme necessidade)
    sleep(1);
}
