package test

import (
	"time"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WS Behavior", func() {

	Context("with two users", func() {

		BeforeEach(func() {
			postStrokeUser1, postStrokeUser1Byte = createPostStroke(username1, []float64{-79.38066843, 43.65483486})
			postStrokeUser2, postStrokeUser2Byte = createPostStroke(username2, []float64{-79.38066843, 43.65483486})
		})

		It("should do match", func() {
			wsConnUser1.WriteMessage(websocket.TextMessage, postStrokeUser1Byte)
			wsConnUser2.WriteMessage(websocket.TextMessage, postStrokeUser2Byte)
			_, resp1, err1 := wsConnUser1.ReadMessage()
			_, resp2, err2 := wsConnUser2.ReadMessage()
			Expect(err1).To(BeNil())
			Expect(string(resp1)).To(BeEquivalentTo(postStrokeUser2.Info))
			Expect(err2).To(BeNil())
			Expect(string(resp2)).To(BeEquivalentTo(postStrokeUser1.Info))
		})

		It("should do match after 1 second", func() {
			wsConnUser1.WriteMessage(websocket.TextMessage, postStrokeUser1Byte)
			time.Sleep(1 * time.Second)
			wsConnUser2.WriteMessage(websocket.TextMessage, postStrokeUser2Byte)
			_, resp1, err1 := wsConnUser1.ReadMessage()
			_, resp2, err2 := wsConnUser2.ReadMessage()
			Expect(err1).To(BeNil())
			Expect(string(resp1)).To(BeEquivalentTo(postStrokeUser2.Info))
			Expect(err2).To(BeNil())
			Expect(string(resp2)).To(BeEquivalentTo(postStrokeUser1.Info))
		})

		It("should do not match", func() {
			wsConnUser1.WriteMessage(websocket.TextMessage, postStrokeUser1Byte)
			time.Sleep(3 * time.Second)
			wsConnUser2.WriteMessage(websocket.TextMessage, postStrokeUser2Byte)
			_, _, err1 := wsConnUser1.ReadMessage()
			_, _, err2 := wsConnUser2.ReadMessage()
			Expect(err1).NotTo(BeNil())
			Expect(err2).NotTo(BeNil())
		})

	})

	Context("with two users far away", func() {

		BeforeEach(func() {
			postStrokeUser1, postStrokeUser1Byte = createPostStroke(username1, []float64{-79.38066843, 43.65483486})
			postStrokeUser2, postStrokeUser2Byte = createPostStroke(username2, []float64{-49.38066843, 43.65483486})
		})

		It("should do not match", func() {
			wsConnUser1.WriteMessage(websocket.TextMessage, postStrokeUser1Byte)
			time.Sleep(3 * time.Second)
			wsConnUser2.WriteMessage(websocket.TextMessage, postStrokeUser2Byte)
			_, _, err1 := wsConnUser1.ReadMessage()
			_, _, err2 := wsConnUser2.ReadMessage()
			Expect(err1).NotTo(BeNil())
			Expect(err2).NotTo(BeNil())
		})

	})

	Context("with three users", func() {

		BeforeEach(func() {
			postStrokeUser1, postStrokeUser1Byte = createPostStroke(username1, []float64{-79.38066843, 43.65483486})
			postStrokeUser2, postStrokeUser2Byte = createPostStroke(username2, []float64{-79.38066843, 43.65483486})
			postStrokeUser3, postStrokeUser3Byte = createPostStroke(username3, []float64{-79.38066843, 43.65483486})
		})

		It("should do match", func() {
			wsConnUser1.WriteMessage(websocket.TextMessage, postStrokeUser1Byte)
			wsConnUser2.WriteMessage(websocket.TextMessage, postStrokeUser2Byte)
			wsConnUser3.WriteMessage(websocket.TextMessage, postStrokeUser3Byte)
			matchOtherTwo(wsConnUser1, postStrokeUser2.Info, postStrokeUser3.Info)
			matchOtherTwo(wsConnUser2, postStrokeUser1.Info, postStrokeUser3.Info)
			matchOtherTwo(wsConnUser3, postStrokeUser2.Info, postStrokeUser1.Info)
		})

		It("a1 and a2 should match, 3 shouldn't match", func() {
			wsConnUser1.WriteMessage(websocket.TextMessage, postStrokeUser1Byte)
			wsConnUser2.WriteMessage(websocket.TextMessage, postStrokeUser2Byte)
			time.Sleep(3 * time.Second)
			wsConnUser3.WriteMessage(websocket.TextMessage, postStrokeUser3Byte)
			_, resp1, err1 := wsConnUser1.ReadMessage()
			_, resp2, err2 := wsConnUser2.ReadMessage()
			Expect(err1).To(BeNil())
			Expect(string(resp1)).To(BeEquivalentTo(postStrokeUser2.Info))
			Expect(err2).To(BeNil())
			Expect(string(resp2)).To(BeEquivalentTo(postStrokeUser1.Info))
			_, _, err3 := wsConnUser3.ReadMessage()
			Expect(err3).NotTo(BeNil())
		})

	})

})
