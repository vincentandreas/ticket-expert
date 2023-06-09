Ticket Expert
Under Construction...

Simulate buy ticket. Using Redis to manage the waiting queue. 

Feature:
- Waiting queue
Website will periodically get their queue status, once other people already finish in order room, the API will notify user, with this response
resp:
```json
{
  "userId":"001",
  "qUniqueCode":"xxx"
}
```

Booking
Resp
```json
{
    "user_id": 1,
    "event_id": 2,
    "qUniqueCode":"xxx",
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
