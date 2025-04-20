import api from "./apiClient";

export interface ProductDetailsType {
    id: number;
    title: string;
    price: number;
    basePrice: number;
    discountPercentage: number;
    description: string;
    isFavorite: boolean;
    averageRating: number;
    images: ProductImage[];
  }

  export interface SimpleProductObject {
    id: number;
    title: string;
    price: number;
    basePrice: number;
    averageRating: number;
    image: string;
    isFavorite: boolean;
  }
  
  export interface ProductImage {
    id: number;
    imageUrl: string;
    sortOrder: number;
  }

  class ProductService {
    async getProductDetails(productId: number): Promise<ProductDetailsType> {
        try {
          const response = await api.get(`/product/details/${productId}`);
          return response.data;
        } catch (error) {
          throw this.handleError(error);
        }
      }
    
      async getAllProducts(): Promise<SimpleProductObject[]> {
        try {
          const response = await api.get("/product/all/details");
          console.log("get All products", response.status)

          return response.data.map((item: any) => ({
            ...item,
            image: item.image.imageUrl
          }));
        } catch (error) {
          throw this.handleError(error);
        }
      }
      
      async addFavorite(productId: number) {
          const response = await api.post(`/favorite/${productId}`);
          return response.data;
      }

      async removeFavorite(productId: number) {
        const response = await api.delete(`/favorite/${productId}`);
        return response.data;
    }

      private handleError(error: any): Error {
        const defaultMessage = "Failed to fetch product data";
        if (error.message) {
          return new Error(error.message || defaultMessage);
        }
        return new Error(defaultMessage);
      }
  }

  export default new ProductService();