package helpers

import (
	"io/ioutil"
	"log"
	"net/http"
	"net"
	"strings"
	"strconv"
	"github.com/gin-gonic/gin"
)

func HttpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return "error", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

func InetNtoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3],bytes[2],bytes[1],bytes[0])
}

// Convert net.IP to int64 ,  http://www.outofmemory.cn
func InetAton(ip string) int64 {
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	bits := strings.Split(ip, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

func Throw403(c *gin.Context, message string){
	c.JSON(http.StatusForbidden, gin.H{
		"code": http.StatusForbidden,
		"message": message,
	})
	c.Abort()
}

func Success(c *gin.Context, data interface{}){
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"status": "success",
		"data": data,
	})
	c.Abort()
}

func Failed(c *gin.Context, status_code int, message string){
	c.JSON(status_code, gin.H{
		"code": status_code,
		"status": "failed",
		"message": message,
	})
	c.Abort()
}