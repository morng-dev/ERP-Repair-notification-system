# ERP Repair Notification System

โปรเจกต์ระบบ ERP สำหรับจัดการงานแจ้งซ่อมและการแจ้งเตือน พัฒนาด้วยภาษา Go (Golang) โดยเน้นการออกแบบ Backend ตามแนวทาง Clean Architecture และ Hexagonal Architecture เพื่อแยกส่วนการทำงานให้ชัดเจนและง่ายต่อการบำรุงรักษา

## ฟีเจอร์หลัก

- REST API สำหรับจัดการข้อมูลและงานแจ้งซ่อม
- ระบบยืนยันตัวตนด้วย JWT และควบคุมสิทธิ์การเข้าถึงด้วย RBAC
- จัดการข้อมูลผ่าน PostgreSQL ด้วย GORM
- ใช้ Redis สำหรับ Cache เพื่อลดการอ่านข้อมูลจากฐานข้อมูลซ้ำ
- ใช้ Apache Kafka สำหรับรับส่ง Event และประมวลผลการแจ้งเตือนแบบ Asynchronous
- ใช้ WebSocket สำหรับส่งการแจ้งเตือนแบบ Real-time
- ใช้ Docker และ Docker Compose สำหรับจัดการสภาพแวดล้อมในการพัฒนา

## เทคโนโลยีที่ใช้

- Go (Golang)
- Fiber
- GORM
- PostgreSQL
- Redis
- Apache Kafka
- WebSocket
- Docker & Docker Compose

## วัตถุประสงค์

เพื่อพัฒนาทักษะ Go Backend Development เรียนรู้การออกแบบระบบที่แยกความรับผิดชอบอย่างชัดเจน และทำความเข้าใจการทำงานร่วมกันของ Database, Cache, Message Broker และ Real-time Communication

## สถานะโปรเจกต์

อยู่ระหว่างการพัฒนาและปรับปรุง เพื่อเพิ่มความเข้าใจด้าน Backend Engineering และสร้างผลงานสำหรับ Portfolio
