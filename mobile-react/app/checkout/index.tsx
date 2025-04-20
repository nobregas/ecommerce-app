import { View, Text, Alert, ScrollView, StyleSheet, TouchableOpacity, ActivityIndicator } from 'react-native'
import React, { useEffect, useState } from 'react'
import { getTotal } from '@/utils/sharedFunctions'
import { Colors } from '@/constants/Colors'
import { Stack, useNavigation, useRouter } from 'expo-router'
import { useHeaderHeight } from '@react-navigation/elements'
import { CartItemType } from '@/types/type'
import { Image } from 'expo-image'
import Animated, { SlideInDown } from 'react-native-reanimated'
import Icon from 'react-native-vector-icons/MaterialCommunityIcons'
import Ionicons from 'react-native-vector-icons/Ionicons'
import { useCartStore } from '@/store/cardBadgeStore'
import { useOrderStore } from '@/store/orderStore'
import { PaymentMethod } from '@/service/orderService'
import cartService from '@/service/cartService'

const CheckoutScreen = () => {
    const navigation = useNavigation()
    const router = useRouter()
    const headerHeight = useHeaderHeight()
    const { cartItems, fetchCartItems, total } = useCartStore()
    const { createOrder, isLoading, error, clearError } = useOrderStore()
    const [selectedPayment, setSelectedPayment] = useState<PaymentMethod | null>(null)

    useEffect(() => {
        fetchCartItems()
    }, [])

    useEffect(() => {
        if (error) {
            Alert.alert("Error", error, [
                { text: "OK", onPress: clearError }
            ]);
        }
    }, [error])

    const handlePlaceOrder = async () => {
        try {
            if (!selectedPayment) {
                Alert.alert("Error", "Please select a payment method.")
                return
            }

            const payload = {
                paymentMethod: selectedPayment,
                paymentId: selectedPayment === PaymentMethod.PIX ? 'pix-' + Date.now() : undefined
            }

            const order = await createOrder(payload)
            
            // Clear the cart after order is placed successfully
            await cartService.clearCart()
            // Refresh cart
            fetchCartItems()
            
            Alert.alert(
                "Order Placed", 
                "Your order has been successfully placed! Order #" + order.id,
                [
                    { 
                        text: "View Details", 
                        onPress: () => router.push(
                            "/orders/index"
                        )
                    },
                    { 
                        text: "Continue Shopping", 
                        onPress: () => router.push("/(tabs)") 
                    }
                ]
            )
        } catch (error) {
            console.error("Error creating order:", error)
            Alert.alert("Error", "Failed to place order. Please try again.")
        }
    }

    const getPaymentIcon = (method: PaymentMethod) => {
        switch (method) {
            case PaymentMethod.CREDIT_CARD:
                return 'credit-card'
            case PaymentMethod.DEBIT_CARD:
                return 'bank-transfer'
            case PaymentMethod.PIX:
                return 'qrcode-scan'
        }
    }

    const getPaymentLabel = (method: PaymentMethod) => {
        switch (method) {
            case PaymentMethod.CREDIT_CARD:
                return 'Credit Card'
            case PaymentMethod.DEBIT_CARD:
                return 'Debit Card'
            case PaymentMethod.PIX:
                return 'PIX'
        }
    }

    if (cartItems.length === 0) {
        return (
            <>
                <Stack.Screen
                    options={{ headerShown: true, headerTransparent: true, headerTitleAlign: 'center', title: 'Checkout' }}
                />
                <View style={[styles.emptyContainer, { marginTop: headerHeight }]}>
                    <Text style={styles.emptyText}>Your cart is empty</Text>
                    <TouchableOpacity 
                        style={styles.continueShoppingBtn}
                        onPress={() => router.push("/(tabs)")}
                    >
                        <Text style={styles.continueShoppingBtnTxt}>Continue Shopping</Text>
                    </TouchableOpacity>
                </View>
            </>
        )
    }

    return (
        <>
            <Stack.Screen
                options={{ headerShown: true, headerTransparent: true, headerTitleAlign: 'center', title: 'Checkout' }}
            />
            <ScrollView
                style={[styles.container, { marginTop: headerHeight }]}
            >
                <View style={styles.section}>
                    <Text style={styles.sectionTitle}>Order Summary</Text>

                    {cartItems.map((item, index) => (
                        <View key={index} style={styles.itemContainer}>
                            <Image
                                source={{ uri: item.productImage }}
                                style={styles.itemImage}
                                contentFit="cover"
                            />
                            <View style={styles.itemInfoWrapper}>
                                <Text style={styles.itemName}>{item.productTitle}</Text>
                                <Text style={styles.itemDetails}>
                                    {item.quantity} x R${item.priceAtAdding.toFixed(2)}
                                </Text>
                            </View>
                        </View>
                    ))}
                </View>

                <View style={styles.section}>
                    <Text style={styles.sectionTitle}>Payment Method</Text>
                    
                    {Object.values(PaymentMethod).map((method) => (
                        <TouchableOpacity
                        key={method}
                        style={[
                            styles.paymentOption,
                            selectedPayment === method && styles.selectedOption
                        ]}
                        onPress={() => setSelectedPayment(method)}
                    >
                        <Icon
                            name={getPaymentIcon(method)}
                            size={24}
                            color={selectedPayment === method ? Colors.primary : Colors.darkGray}
                        />
                        <Text style={styles.paymentText}>
                            {getPaymentLabel(method)}
                        </Text>
                        {selectedPayment === method && (
                            <Ionicons
                                name="checkmark-circle"
                                size={24}
                                color={Colors.primary}
                                style={styles.checkmark}
                            />
                        )}
                    </TouchableOpacity>
                    ))}
                
                </View>

            </ScrollView>
            <Animated.View style={styles.footer} entering={SlideInDown.duration(500)}>
                <View style={styles.priceInfoWrapper}>
                    <Text style={styles.totalTxt}>Total: R${total.toFixed(2)}</Text>
                </View>
                <TouchableOpacity
                    style={styles.checkoutBtn}
                    onPress={handlePlaceOrder}
                    disabled={isLoading}
                >
                    {isLoading ? (
                        <ActivityIndicator color={Colors.white} />
                    ) : (
                        <Text style={styles.checkoutBtnTxt}>Place Order</Text>
                    )}
                </TouchableOpacity>
            </Animated.View>
        </>
    )
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: Colors.background,
        padding: 20,
        marginBottom: 80
    },
    section: {
        backgroundColor: Colors.white,
        borderRadius: 10,
        padding: 15,
        marginBottom: 20,
        shadowColor: Colors.darkGray,
        shadowOffset: { width: 0, height: 2 },
        shadowOpacity: 0.1,
        shadowRadius: 4,
        elevation: 3,
    },
    sectionTitle: {
        fontSize: 18,
        fontWeight: '600',
        color: Colors.primary,
        marginBottom: 15,
        borderBottomWidth: 1,
        borderBottomColor: Colors.extraLightGray,
        paddingBottom: 10,
    },
    itemContainer: {
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
    itemImage: {
        width: 60,
        height: 60,
        borderRadius: 8,
        marginRight: 15,
        flexShrink: 0,
    },
    itemInfo: {
        flex: 1,
        flexShrink: 1,
        marginRight: 10,
    },
    itemName: {
        fontSize: 16,
        fontWeight: '500',
        color: Colors.black,
        minWidth: 100,
        maxWidth: "80%",
    },
    itemDetails: {
        fontSize: 14,
        color: Colors.darkGray,
        fontWeight: '400',
        flexWrap: 'wrap',
    },
    paymentOption: {
        flexDirection: 'row',
        alignItems: 'center',
        padding: 15,
        marginBottom: 10,
        borderWidth: 1,
        borderColor: Colors.lightGray,
        borderRadius: 8,
    },
    selectedOption: {
        borderColor: Colors.primary,
        borderWidth: 2,
        backgroundColor: Colors.extraLightGray,
    },
    paymentText: {
        fontSize: 16,
        marginLeft: 15,
        color: Colors.darkGray,
        flex: 1,
    },
    checkmark: {
        marginLeft: 'auto',
    },
    paymentNote: {
        backgroundColor: Colors.transparentWhite,
        borderRadius: 8,
        padding: 15,
        marginTop: 10,
    },
    noteText: {
        color: Colors.gray,
        fontSize: 14,
        lineHeight: 20,
    },
    footer: {
        position: 'absolute',
        bottom: 0,
        left: 0,
        right: 0,
        width: '100%',
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginTop: 20,
        paddingVertical: 15,
        paddingHorizontal: 20,
        backgroundColor: Colors.white,
        borderTopWidth: 1,
        borderTopColor: Colors.extraLightGray,
        shadowColor: Colors.darkGray,
        shadowOffset: { width: 0, height: -2 },
        shadowOpacity: 0.1,
        shadowRadius: 4,
        elevation: 5,
    },
    priceInfoWrapper: {
        flex: 1,
        justifyContent: 'center'
    },
    totalTxt: {
        fontSize: 16,
        fontWeight: '500',
        color: Colors.black
    },
    checkoutBtn: {
        flex: 1,
        backgroundColor: Colors.primary,
        height: 46,
        justifyContent: 'center',
        alignItems: 'center',
        borderRadius: 8,
    },
    checkoutBtnTxt: {
        fontSize: 16,
        fontWeight: '600',
        color: Colors.white
    },
    emptyContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: Colors.background,
        padding: 20,
    },
    emptyText: {
        fontSize: 18,
        color: Colors.gray,
        fontWeight: '500',
        marginBottom: 20,
    },
    continueShoppingBtn: {
        backgroundColor: Colors.primary,
        paddingVertical: 12,
        paddingHorizontal: 24,
        borderRadius: 8,
    },
    continueShoppingBtnTxt: {
        color: Colors.white,
        fontWeight: '600',
        fontSize: 16,
    }
});

export default CheckoutScreen;