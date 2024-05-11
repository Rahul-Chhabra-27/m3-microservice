package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"rahulchhabra.io/config"
	"rahulchhabra.io/model"
	orderproto "rahulchhabra.io/proto/order"
)

type OrderService struct {
	orderproto.UnimplementedOrderServiceServer
}

func (*OrderService) GetAllOrder(ctx context.Context, req *orderproto.OrderRequest) (*orderproto.OrderResponse, error) {
	// get all the orders from the database
	orders, err := model.OrderCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, config.ErrorMessage("Could not fetch orders", codes.Internal)
	}
	var orderList []*orderproto.Order
	for orders.Next(context.Background()) {
		var order model.Order
		orders.Decode(&order)
		orderList = append(orderList, &orderproto.Order{
			Id:       order.Id.Hex(),
			Name:     order.Name,
			Price:    order.Price,
			Quantity: order.Quantity,
		})
	}
	return &orderproto.OrderResponse{Orders: orderList}, nil
}
func InsertOrders() {
	check := model.OrderCollection.FindOne(context.Background(), model.Order{Name: "Pizza"})

	if check.Err() == nil {
		return
	}
	// Create a list of orders
	orders := []model.Order{
		{
			Name:     "Pizza",
			Price:    100,
			Quantity: 1,
		},
		{
			Name:     "Burger",
			Price:    50,
			Quantity: 1,
		},
		{
			Name:     "Pasta",
			Price:    80,
			Quantity: 1,
		},
		{
			Name:     "Sandwich",
			Price:    30,
			Quantity: 1,
		},
		{
			Name:     "Coke",
			Price:    20,
			Quantity: 1,
		},
		{
			Name:     "Pepsi",
			Price:    25,
			Quantity: 1,
		},
	}

	// Insert all the orders into the database
	for _, order := range orders {
		_, err := model.OrderCollection.InsertOne(context.Background(), order)
		if err != nil {
			log.Fatal("Failed to insert order: ", err)
			return
		}
	}
	fmt.Println("Orders inserted successfully")
}

// Responsible for starting the server
func startServer() {
	godotenv.Load()
	// Log a message
	fmt.Println("Starting server...")
	// Initialize the gotenv file..
	godotenv.Load()

	// Create a new context
	ctx := context.TODO()

	// Connect to the MongoDB database
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	// Set the global variable to the collection
	model.OrderCollection = db.Database("testdb").Collection("dummyorders")
	InsertOrders()
	// Check for errors
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s", err)
	}

	listner, err := net.Listen("tcp", "localhost:50052")
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
	fmt.Println("Database connected Successfully")

	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	orderproto.RegisterOrderServiceServer(grpcServer, &OrderService{})

	go func() {
		if err := grpcServer.Serve(listner); err != nil {
			log.Fatalf("Failed to serve: %s", err)
		}
	}()
	// Create a new gRPC-Gateway server (gateway).
	connection, err := grpc.DialContext(
		context.Background(),
		"localhost:50052",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	// Create a new gRPC-Gateway mux (gateway).
	gwmux := runtime.NewServeMux()

	orderproto.RegisterOrderServiceHandler(context.Background(), gwmux, connection)
	gwServer := &http.Server{
		Addr:    ":8091",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8091")
	log.Fatalln(gwServer.ListenAndServe())
}
func main() {
	startServer()
}
