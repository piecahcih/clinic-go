package appointment

import (
	"context"
	"log"
	"time"
)

func RunNoShowWorker(ctx context.Context, repo Repository, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	sem := make(chan struct{}, 1)
	sem <- struct{}{} // push struct in sem, now sem full

	for {
		select {
		case <-ctx.Done():
			log.Println("no-show worker: shutting down")
			return
		case <-ticker.C:
			select {
			case <-sem: //เอาของออกจากsem ให้มันfree
				go func() {
					defer func() { sem <- struct{}{} }() // release when done
					n, err := repo.MarkNoShows(ctx)
					if err != nil {
						log.Printf("mark no-shows: %v", err)
						return
					}
					if n > 0 {
						log.Printf("marked %d appointment(s) as no_show", n)
					}
				}()
			default:
				log.Println("no-show worker: previous sweep still running, skipping tick")
			}
		}
	}
}

//logic นี้ คือ จะเริ่มทำก็ต่อเมื่อมีของในchannel(1 เต็ม) ในขณะที่ทำเอาของออกทำให้
// อันอื่นทำไม่ได้จนกว่าฟังก์ชั่นนั้นจะจบแล้วเติมของกลับตามที่เซตในdefer
// sweep รอบก่อนหน้ายังไม่เสร็จ อย่าเพิ่งเริ่มรอบใหม่"
