import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '30s', target: 1000 },
        { duration: '1m', target: 2000 },
        { duration: '30s', target: 0 },
    ],
};

let globalAmount = 1;

export default function () {
    const url = 'http://localhost:32300/auctions/1/bids';

    const currentAmount = globalAmount + (__VU * 10) + __ITER;

    const payload = JSON.stringify({
        account_email: "juliocruz@discente.ufg.br",        amount: currentAmount,
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