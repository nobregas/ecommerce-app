
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native'
import React from 'react'
import { Colors } from '@/constants/Colors'
import { Ionicons } from '@expo/vector-icons'
import { Image } from "expo-image"
import { Link } from 'expo-router'
import { useCartStore } from '@/store/cardBadgeStore'
import { CartItemType } from '@/service/cartService'

type Props = {
    item: CartItemType
}

const CartItem = ({ item }: Props) => {
    const { addToCart, removeOneFromCart, removeItemFromCart } = useCartStore();
    
    const handleIncrement = async () => {
        try {
            await addToCart(item.productId);
        } catch (error) {
            alert("Error incrementing item");
        }
    }

    const handleDecrement = async () => {
        try {
            await removeOneFromCart(item.productId);
        } catch (error) {
            alert("Error decrementing item");
        }
    }

    const handleRemoveItem = async () => {
        try {
            await removeItemFromCart(item.productId);
        } catch (error) {
            alert("Error removing item");
        }
    }


    return (
        <Link href={{
            pathname: "/product-details/[id]",
            params: { id: item.productId }
        }} asChild>
            <TouchableOpacity>
                <View style={styles.container}>
                    <Image 
                        source={{ uri: item.productImage }} 
                        style={styles.itemImg} 
                        contentFit="cover"
                    />
                    <View style={styles.itemInfoWrapper}>
                        <Text style={styles.itemTxt} numberOfLines={1}>{item.productTitle}</Text>
                        <Text style={styles.itemTxt}>R${item.priceAtAdding.toFixed(2)}</Text>
                        <View style={styles.itemControlWrapper}>
                            <TouchableOpacity 
                                onPress={(e) => {
                                    e.preventDefault();
                                    handleRemoveItem();
                                }}
                            >
                                <Ionicons name="trash-outline" size={20} color={Colors.red} />
                            </TouchableOpacity>
                            <View style={styles.quantityControlWrapper}>
                                <TouchableOpacity 
                                    style={styles.quantityControl} 
                                    onPress={(e) => {
                                        e.preventDefault();
                                        handleDecrement();
                                    }}
                                >
                                    <Ionicons name="remove-outline" size={20} color={Colors.black} />
                                </TouchableOpacity>
                                <Text style={styles.quantity}>{item.quantity}</Text>
                                <TouchableOpacity 
                                    style={styles.quantityControl} 
                                    onPress={(e) => {
                                        e.preventDefault();
                                        handleIncrement();
                                    }}
                                >
                                    <Ionicons name="add-outline" size={20} color={Colors.black} />
                                </TouchableOpacity>
                            </View>
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