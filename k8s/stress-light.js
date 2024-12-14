import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '30s', target: 10 },  // Aumenta de 0 para 10 usuários simultâneos em 30 segundos
        { duration: '1m', target: 50 },   // Aumenta para 50 usuários em 1 minuto
        { duration: '30s', target: 0 },   // Diminui de 50 para 0 usuários em 30 segundos
    ],
};

let globalAmount = 1;

export default function () {
    const url = 'http://localhost:32300/auctions/3/bids';

    const currentAmount = globalAmount + (__VU * 10) + __ITER;

    const payload = JSON.stringify({
        account_email: "juliocruz@discente.ufg.br",
        amount: currentAmount,
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const response = http.post(url, payload, params);

    check(response, {
        'status is 202': (r) => r.status === 202,
    });

    console.log(`Virtual User: ${__VU}, Iteração: ${__ITER}, Amount: ${currentAmount}`);

    sleep(1);
}