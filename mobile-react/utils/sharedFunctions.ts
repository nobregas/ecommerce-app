import { CartItemType } from "@/types/type";

export const getTotal = (cartItems: CartItemType[]) => {
    let total = 0;
    for (let i = 0; i < cartItems.length; i++) {
        total += cartItems[i].price * cartItems[i].quantity
    }
    return total.toFixed(2)
}
export const isValidEmail = (email: string) => {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
};

export const isValidCPF = (cpf: string) => {
  cpf = cpf.replace(/[^\d]+/g, '');
  if (cpf.length !== 11) return false;
  
  return /^\d{3}\.\d{3}\.\d{3}-\d{2}$/.test(cpf) || cpf.length === 11;
};

export const formatCPF = (value: string) => {
    return value
        .replace(/\D/g, '')
        .replace(/(\d{3})(\d)/, '$1.$2')
        .replace(/(\d{3})(\d)/, '$1.$2')
        .replace(/(\d{3})(\d)/, '$1-$2')
        .substring(0, 14);
};

export const formatDateTime = (dateTimeStr: string): string => {
  try {
    const date = new Date(dateTimeStr);
    
    if (isNaN(date.getTime())) {
      return dateTimeStr; 
    }
    
    const day = date.getDate().toString().padStart(2, '0');
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const year = date.getFullYear();
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    
    return `${day}/${month}/${year} ${hours}:${minutes}`;
  } catch (error) {
    return dateTimeStr; 
  }
}