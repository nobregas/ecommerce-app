import { Platform } from "react-native";
import axios from "axios";
import { GEMINI_API_KEY, GEMINI_API_URL } from '@env';
import productService, { ProductDetailsType, SimpleProductObject } from "./productService";
import categoryService, { Category } from "./categoryService";
import userService, { UserDTO } from "./userService";
import cartService, { CartItemType } from "./cartService";
import orderService, { OrderHistory, OrderWithItems } from "./orderService";

export interface ChatMessageType {
  id: string;
  content: string;
  role: "user" | "assistant";
  timestamp: Date;
}

export interface ChatSessionType {
  id: string;
  messages: ChatMessageType[];
  createdAt: Date;
}

export interface IAIChatService {
  sendMessage(message: string, sessionId?: string): Promise<ChatMessageType>;
}

export class AIChatService implements IAIChatService {
  private apiKey: string;
  private baseUrl: string;
  private systemPrompt: string;
  
  constructor() {
    this.apiKey = GEMINIKEY;
    this.baseUrl = GEMINIMODELURL;
    
    this.systemPrompt = `
You are a helpful customer support assistant for an e-commerce app called ShopX.
The app sells fashion items, electronics, home goods, and accessories.

Key information about ShopX:
- Standard delivery takes 3-5 business days
- Express delivery takes 1-2 business days
- Return policy: 30 days for unworn/unused items with original packaging
- Payment methods: Credit/debit cards, PayPal, Apple Pay, Google Pay
- For order issues, customers should provide their order number
- ShopX operates in the United States, Canada, and selected European countries
- Customer service phone number: (800) 555-1234
- Customer service email: support@shopx.com

Please provide helpful, concise, and friendly support to the customer.
Suggest specific ShopX policies and services when relevant to the customer's query.
If you don't know the answer to a specific product question, politely ask for more details or suggest contacting customer service.
`;
  }

  async sendMessage(message: string, sessionId?: string): Promise<any> {
    try {
      // Get current context data
      let contextData = await this.gatherContextData();
    

      const payload = {
        contents: [
          {
            role: "user",
            parts: [{ text: `${this.systemPrompt}\n\n${this.formatContextData(contextData)}\n\nUser Query: ${message}` }]
          },
        ],
          systemInstruction: {
            parts: [{ text: this.systemPrompt }]
          },
          generationConfig: {
            temperature: 0.7,
            maxOutputTokens: 800,
          }
      };

      console.log("Sending request to Gemini API");
      
      const response = await axios.post(
        `${this.baseUrl}?key=${this.apiKey}`,
        payload,
        {
          headers: {
            'Content-Type': 'application/json'
          }
        }
      );

      console.log("Received response from Gemini API");
      
      if (response.data && 
          response.data.candidates && 
          response.data.candidates[0] && 
          response.data.candidates[0].content && 
          response.data.candidates[0].content.parts && 
          response.data.candidates[0].content.parts[0]) {
        
        const content = response.data.candidates[0].content.parts[0].text;
        
        return {
          id: Date.now().toString(),
          content,
          role: "assistant",
          timestamp: new Date(),
        };
      } else {
        console.error("Unexpected response format from Gemini API:", response.data);
      }
    } catch (error) {
      console.error("Error sending message to Gemini:", error);
      
      if (axios.isAxiosError(error)) {
        console.error("API Error details:", {
          status: error.response?.status,
          data: error.response?.data,
          headers: error.response?.headers
        });

        if (error.response?.status === 404) {
          console.error("404 error received. Check if the API endpoint is correct");
        }
      }
      
      const contextData = await this.gatherContextData();
    }
  }

  private async gatherContextData(): Promise<{
    user?: UserDTO;
    products?: SimpleProductObject[];
    categories?: Category[];
    cartItems?: CartItemType[];
    orders?: OrderHistory[];
  }> {
    const contextData: {
      user?: UserDTO;
      products?: SimpleProductObject[];
      categories?: Category[];
      cartItems?: CartItemType[];
      orders?: OrderHistory[];
    } = {};

    try {
      // Get current user data
      contextData.user = await userService.getCurrentUser();
    } catch (error) {
      console.log("Could not get user data:", error);
    }

    try {
      // Get products
      contextData.products = await productService.getAllProducts();
    } catch (error) {
      console.log("Could not get products:", error);
    }

    try {
      // Get categories
      contextData.categories = await categoryService.getAllCategories();
    } catch (error) {
      console.log("Could not get categories:", error);
    }

    try {
      // Get cart items
      contextData.cartItems = await cartService.getCartItems();
    } catch (error) {
      console.log("Could not get cart items:", error);
    }

    try {
      // Get orders
      contextData.orders = await orderService.getOrders();
    } catch (error) {
      console.log("Could not get orders:", error);
    }

    return contextData;
  }

  private formatContextData(contextData: any): string {
    let formattedData = "\n\nCURRENT USER CONTEXT:\n";
    
    if (contextData.user) {
      formattedData += `\nUSER:\n`;
      formattedData += `Name: ${contextData.user.fullName}\n`;
      formattedData += `Email: ${contextData.user.email}\n`;
      formattedData += `Customer since: ${contextData.user.createdAt}\n`;
    }
    
    if (contextData.cartItems && contextData.cartItems.length > 0) {
      formattedData += `\nSHOPPING CART (${contextData.cartItems.length} items):\n`;
      contextData.cartItems.forEach((item: CartItemType) => {
        formattedData += `- ${item.quantity}x ${item.productTitle} - $${item.priceAtAdding.toFixed(2)}\n`;
      });
    } else {
      formattedData += `\nSHOPPING CART: Empty\n`;
    }
    
    if (contextData.orders && contextData.orders.length > 0) {
      formattedData += `\nRECENT ORDERS:\n`;
      contextData.orders.slice(0, 3).forEach((order: OrderHistory) => {
        formattedData += `- Order #${order.id}: $${order.totalAmount.toFixed(2)} - Status: ${order.status} - Date: ${order.createdAt}\n`;
      });
    }
    
    if (contextData.categories) {
      formattedData += `\nAVAILABLE CATEGORIES:\n`;
      contextData.categories.forEach((category: Category) => {
        formattedData += `- ${category.name}\n`;
      });
    }
    
    if (contextData.products) {
      formattedData += `\nPOPULAR PRODUCTS (first 5):\n`;
      contextData.products.slice(0, 5).forEach((product: SimpleProductObject) => {
        const discount = product.basePrice > product.price 
          ? ` (${Math.round((1 - product.price/product.basePrice) * 100)}% discount)`
          : '';
        formattedData += `- ${product.title} - $${product.price.toFixed(2)}${discount} - Rating: ${product.averageRating}/5\n`;
      });
    }

    return formattedData;
  }

}

export const aiChatService = new AIChatService(); 
