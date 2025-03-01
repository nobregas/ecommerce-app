import { CartItemType } from "@/types/type";

export const getTotal = (cartItems: CartItemType[]) => {
    let total = 0;
    for (let i = 0; i < cartItems.length; i++) {
        total += cartItems[i].price * cartItems[i].quantity
    }
    return total.toFixed(2)
}