
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native'
import React from 'react'
import { CartItemType } from '@/mobile/types/type'
import { Colors } from '@/mobile/constants/Colors'
import { Ionicons } from '@expo/vector-icons'
import { Image } from "expo-image"
import { Link } from 'expo-router'
import { useCartStore } from '@/store/cardBadgeStore'


type Props = {
    item: CartItemType
    updateQuantity: (id: number, newQuantity: number) => void
    removeItem: (id: number) => void
}


const CartItem = ({ item, updateQuantity, removeItem }: Props) => {
    const { fetchCartCount } = useCartStore();

    const handleIncrement = async () => {
        await updateQuantity(item.id, item.quantity + 1)
    }

    const handleDecrement = async () => {
        if (item.quantity > 1) {
            await updateQuantity(item.id, item.quantity - 1)
        }
    }

    const handleRemoveItem = async () => {
        await removeItem(item.id)
        fetchCartCount()
    }

    return (
        <Link href={{
            pathname: "/product-details/[id]",
            params: { id: item.id, productType: item.productType }
        }} asChild>
            <TouchableOpacity>
            <View style={styles.container}>
                <Image source={{ uri: item.image }} style={styles.itemImg} />
                <View style={styles.itemInfoWrapper}>
                    <Text style={styles.itemTxt}>{item.title}</Text>
                    <Text style={styles.itemTxt}>R${item.price}</Text>
                    <View style={styles.itemControlWrapper}>
                        <TouchableOpacity onPress={handleRemoveItem}>
                            <Ionicons name="trash-outline" size={20} color={Colors.red} />
                        </TouchableOpacity>
                        <View style={styles.quantityControlWrapper}>
                            <TouchableOpacity style={styles.quantityControl} onPress={handleDecrement}>
                                <Ionicons name="remove-outline" size={20} color={Colors.black} />
                            </TouchableOpacity>
                            <Text style={styles.quantity}>{item.quantity}</Text>
                            <TouchableOpacity style={styles.quantityControl} onPress={handleIncrement}>
                                <Ionicons name="add-outline" size={20} color={Colors.black} />
                            </TouchableOpacity>
                        </View>
                        <TouchableOpacity>
                            <Ionicons name="heart-outline" size={20} color={Colors.black} />
                        </TouchableOpacity>
                    </View>
                </View>
            </View>
            </TouchableOpacity>
        </Link>
    )
}

export default CartItem

const styles = StyleSheet.create({
    container: {
        flexDirection: 'row',
        alignItems: 'center',
        padding: 10,
        marginBottom: 10,
        borderWidth: StyleSheet.hairlineWidth,
        borderColor: Colors.lightGray,
        borderRadius: 5,
        backgroundColor: Colors.white,
        width: "100%"
    },
    itemInfoWrapper: {
        flex: 1,
        alignSelf: "flex-start",
        gap: 10
    },
    itemImg: {
        width: 100,
        height: 100,
        borderRadius: 5,
        marginRight: 10
    },
    itemTxt: {
        fontSize: 16,
        fontWeight: '500',
        color: Colors.black,
        minWidth: 100,
        maxWidth: "80%",
    },
    itemControlWrapper: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        gap: 10,
    },
    quantityControlWrapper: {
        flexDirection: 'row',
        alignItems: 'center',
        gap: 15,
    },
    quantityControl: {
        padding: 5,
        borderWidth: 1,
        borderColor: Colors.lightGray,
        borderRadius: 5,
    },
    quantity: {
        fontSize: 16,
        fontWeight: '500',
        color: Colors.black
    }
});