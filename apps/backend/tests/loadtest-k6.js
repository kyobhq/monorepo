import ws from 'k6/ws'
// import http from 'k6/http'

export const options = {
  scenarios: {
    websocket_connections: {
      executor: 'ramping-vus',
      startVUs: 1,
      stages: [
        { duration: '30s', target: 2500 },
        { duration: '10s', target: 2500 },
        { duration: '30s', target: 0 },
      ],
      exec: 'websocketFlow',
    },
    // http_messages: {
    //   executor: 'ramping-arrival-rate',
    //   startRate: 5,
    //   timeUnit: '1s',
    //   preAllocatedVUs: 10,
    //   maxVUs: 50,
    //   stages: [
    //     { duration: '30s', target: 10 },   // Gentle ramp
    //     { duration: '60s', target: 50 },   // Still manageable
    //     { duration: '30s', target: 10 },   // Scale down
    //   ],
    //   exec: 'httpFlow',
    // },
  },
  noConnectionReuse: true,
  userAgent: 'k6-websocket-test',
  throw: false,
}

function randomId(min, max) {
  return Math.floor(Math.random() * (max - min + 1) + min).toString()
}

export function websocketFlow() {
  const userId = randomId(1, 10000000)
  const url = `ws://localhost:8080/${userId}`

  const res = ws.connect(url, { tags: { name: 'websocket_connection' } }, function (socket) {
    socket.on('open', function open() {
      // console.log(`[${userId}] WebSocket OPENED successfully`)
    })

    socket.on('message', function (message) {
      // console.log(`[${userId}] Received message`)
    })

    socket.on('close', function (code, reason) {
      // console.log(`[${userId}] WebSocket CLOSED - Code: ${code}, Reason: ${reason}`)
    })

    socket.on('error', function (e) {
      console.error(`[${userId}] WebSocket ERROR: ${e.error()}`)
    })
  })

  console.log(`[${userId}] ws.connect returned:`, res)
  console.log(`[${userId}] websocketFlow ending`)
}

// export function httpFlow() {
//   for (let i = 0; i < 100; i++) {
//     const userId = randomId(1, 10000000)
//     const response = http.post(
//       `http://localhost:8080/${userId}/chat`,
//       JSON.stringify({
//         server_id: "global",
//         author_id: userId,
//         content: `Load test message ${i} from ${userId}`
//       }),
//       {
//         headers: {
//           'Content-Type': 'application/json',
//         },
//       }
//     )
//
//     check(response, {
//       'status is 200': (r) => r.status === 200,
//     })
//
//     sleep(2)
//   }
// }
