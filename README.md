Ticket Expert
Under Construction...


Save user in order room
Website will periodically get their queue status, once other people already finish in order room, the API will notify user, with this response
resp:
```json
{
  "userId":"001",
  "queueUniqueCode":"xxx"
}
```

Booking
Resp
```json
{
    "user_id": 1,
    "event_id": 2,
    "queueUniqueCode":"xxx",
    "booking_status": "active",
    "booking_details": [
        {
            "price": "15007",
            "qty": 1,
            "event_detail_id": 3
        }
    ]
}
```
