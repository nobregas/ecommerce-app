export interface ProductType {
  id: number;
  title: string;
  price: number;
  description: string;
  images: string[];
  category: Category;
}

interface Category {
  id: number;
  name: string;
  image: string;
}

export interface CategoryType {
  id: number;
  name: string;
  image: string;
}

export interface CartItemType {
  id: number;
  title: string;
  price: number;
  quantity: number;
  image: string;
  productType: ProductStrType
}

export interface NotificationType {
  id: number;
  title: string;
  message: string;
  timestamp: string;
}

export interface UserType {
  id: number;
  name: string;
  email: string;
  cpf: string;
}

export interface RegisterType {
  name: string;
  email: string;
  password: string;
  confirmPassword: string;
}

export interface LoginType {
  email: string;
  password: string;
}

export type ProductStrType = "regular" | "sale";