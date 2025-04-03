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
    this.apiKey = "AIzaSyDJWUqdPHosEDZ3jT5r2chrZPd0sRPe2aM";
    this.baseUrl = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent";
    
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

  async sendMessage(message: string, sessionId?: string): Promise<ChatMessageType> {
    try {
      // Get current context data
      let contextData = await this.gatherContextData();
      
      if (__DEV__ || !this.apiKey) {
        console.log("Using mock service: API key is missing");
        return this.mockSendMessage(message, contextData);
      }

      const payload = {
        contents: [
          {
            parts: [{ text: this.systemPrompt + this.formatContextData(contextData) }],
            role: "system"
          },
          {
            parts: [{ text: message }],
            role: "user"
          }
        ],
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
        return this.mockSendMessage(message, contextData);
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
      return this.mockSendMessage(message, contextData);
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

  private async mockSendMessage(message: string, contextData: any): Promise<ChatMessageType> {
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    let responseContent = "";
    
    if (message.toLowerCase().includes("hello") || message.toLowerCase().includes("hi") || 
        message.toLowerCase().includes("olá") || message.toLowerCase().includes("oi")) {
      const userName = contextData.user ? contextData.user.fullName : "to ShopX";
      responseContent = `Hello ${userName}! Welcome to ShopX support. How can I help with your shopping today?`;
    } else if (message.toLowerCase().includes("entrega") || message.toLowerCase().includes("envio") || 
               message.toLowerCase().includes("shipping") || message.toLowerCase().includes("delivery")) {
      responseContent = "At ShopX, we offer standard delivery (3-5 business days) and express delivery (1-2 business days). For international orders, shipping times may vary. You can track your order from your account dashboard once it's shipped.";
    } else if (message.toLowerCase().includes("preço") || message.toLowerCase().includes("custo") || 
               message.toLowerCase().includes("price") || message.toLowerCase().includes("cost")) {
      if (contextData.products && contextData.products.length > 0) {
        responseContent = `Our pricing is competitive and clearly marked on each product page. For example, some of our popular products include: ${contextData.products.slice(0, 3).map((p: SimpleProductObject) => `${p.title} for $${p.price.toFixed(2)}`).join(', ')}. If you're looking for the best deals, check out our 'Sale' section or sign up for our newsletter to receive exclusive discounts.`;
      } else {
        responseContent = "Our pricing is competitive and clearly marked on each product page. If you're looking for the best deals, check out our 'Sale' section or sign up for our newsletter to receive exclusive discounts.";
      }
    } else if (message.toLowerCase().includes("devolução") || message.toLowerCase().includes("reembolso") || 
               message.toLowerCase().includes("troca") || message.toLowerCase().includes("return") || 
               message.toLowerCase().includes("refund") || message.toLowerCase().includes("exchange")) {
      responseContent = "ShopX offers a 30-day return policy for unworn/unused items in their original packaging. To initiate a return, go to your order history in your account and select 'Return Items'. Once we receive your return, refunds typically process within 5-7 business days.";
    } else if (message.toLowerCase().includes("pagamento") || message.toLowerCase().includes("cartão") || 
               message.toLowerCase().includes("crédito") || message.toLowerCase().includes("payment") || 
               message.toLowerCase().includes("card") || message.toLowerCase().includes("credit")) {
      responseContent = "ShopX accepts various payment methods including credit/debit cards, PayPal, Apple Pay, and Google Pay. All transactions are secure and encrypted. If you're experiencing payment issues, please try a different payment method or contact your bank.";
    } else if (message.toLowerCase().includes("pedido") || message.toLowerCase().includes("compra") || 
               message.toLowerCase().includes("order") || message.toLowerCase().includes("purchase")) {
      if (contextData.orders && contextData.orders.length > 0) {
        const recentOrder = contextData.orders[0];
        responseContent = `I'd be happy to help with your order. I see that you have a recent order (#${recentOrder.id}) for $${recentOrder.totalAmount.toFixed(2)} with status ${recentOrder.status} from ${recentOrder.createdAt}. Can I help with any specific information about this order?`;
      } else {
        responseContent = "I'd be happy to help with your order. To provide specific information, I'll need your order number which can be found in your confirmation email. With that, I can check the status and provide further assistance.";
      }
    } else if (message.toLowerCase().includes("tamanho") || message.toLowerCase().includes("size") || 
               message.toLowerCase().includes("fit")) {
      responseContent = "ShopX provides detailed size guides on each product page. For clothing items, we recommend checking the specific measurements listed. If you're between sizes, reading customer reviews can often help determine whether to size up or down for a particular item.";
    } else if (message.toLowerCase().includes("desconto") || message.toLowerCase().includes("cupom") || 
               message.toLowerCase().includes("promoção") || message.toLowerCase().includes("discount") || 
               message.toLowerCase().includes("coupon") || message.toLowerCase().includes("promo")) {
      responseContent = "We regularly offer promotions and discounts! The best way to stay updated on our latest deals is to sign up for our newsletter. Additionally, first-time shoppers can get 10% off by subscribing to our emails. Check our homepage banner for any current promotions.";
    } else if (message.toLowerCase().includes("carrinho") || message.toLowerCase().includes("cart")) {
      if (contextData.cartItems && contextData.cartItems.length > 0) {
        const total = contextData.cartItems.reduce((sum: number, item: CartItemType) => sum + (item.priceAtAdding * item.quantity), 0);
        responseContent = `Your cart currently contains ${contextData.cartItems.length} item(s) totaling $${total.toFixed(2)}. Items include: ${contextData.cartItems.map((item: CartItemType) => `${item.quantity}x ${item.productTitle}`).join(', ')}. Can I help you complete your purchase or do you have any questions about these products?`;
      } else {
        responseContent = "Your cart is currently empty. How about exploring some of our popular categories to find products that might interest you?";
      }
    } else if (message.toLowerCase().includes("categoria") || message.toLowerCase().includes("category")) {
      if (contextData.categories && contextData.categories.length > 0) {
        responseContent = `ShopX offers several categories for you to explore, including: ${contextData.categories.map((cat: Category) => cat.name).join(', ')}. Which category are you interested in?`;
      } else {
        responseContent = "ShopX offers various categories including fashion, electronics, home goods, and accessories. You can browse all categories from the main menu of our app.";
      }
    } else {
      const userName = contextData.user ? `, ${contextData.user.fullName}` : "";
      responseContent = `Thank you for contacting ShopX customer support${userName}. I'm here to help with any questions about our products, orders, shipping, returns, or payments. For immediate assistance with specific order issues, you can also reach us at (800) 555-1234 or support@shopx.com. Please provide additional details about how I can assist you.`;
    }
    
    return {
      id: Date.now().toString(),
      content: responseContent,
      role: "assistant",
      timestamp: new Date(),
    };
  }
}

export const aiChatService = new AIChatService(); 