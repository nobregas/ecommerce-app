import { View, Text, Alert, ScrollView, StyleSheet, TouchableOpacity } from 'react-native'
import React, { useEffect, useState } from 'react'
import { getTotal } from '@/utils/sharedFunctions'
import { Colors } from '@/constants/Colors'
import { Stack, useNavigation } from 'expo-router'
import { useHeaderHeight } from '@react-navigation/elements'
import { CartItemType } from '@/types/type'
import { Image } from 'expo-image'
import Animated, { SlideInDown } from 'react-native-reanimated'
import Icon from 'react-native-vector-icons/MaterialCommunityIcons'
import Ionicons from 'react-native-vector-icons/Ionicons'

enum PaymentMethod {
    CREDIT_CARD = 'CREDIT_CARD',
    DEBIT_CARD = 'DEBIT_CARD',
    PIX = 'PIX'
}

const CheckoutScreen = () => {
    const navigation = useNavigation()
    const headerHeight = useHeaderHeight()
    const [cartItems, setCartItems] = useState<CartItemType[]>([])
    const [selectedPayment, setSelectedPayment] = useState<PaymentMethod | null>(null)



    useEffect(() => {
        fetchCartItems()
    }, [])

    const fetchCartItems = async () => {
        
    }


    const handlePlaceOrder = async () => {
        try {
            if (!selectedPayment) {
                Alert.alert("Error", "Please select a payment method.")
                return
            }

            // clear the cart after placing the order

            // logic to place the order
            Alert.alert("Order Placed", "Your order has been successfully placed.");

            navigation.goBack()

        } catch (error) {
            Alert.alert("Error", "Failed to place order. Please try again.");
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

    return (
        <>
            <Stack.Screen
                options={{ headerShown: true, headerTransparent: true, headerTitleAlign: 'center', title: 'checkout' }}
            />
            <ScrollView
                style={[styles.container, { marginTop: headerHeight }]}
            >
                <View style={styles.section}>
                    <Text style={styles.sectionTitle}>Order Summary</Text>

                    {cartItems.map((item, index) => (
                        <View key={index} style={styles.itemContainer}>
                            <Image
                                source={{ uri: item.image }}
                                style={styles.itemImage}
                            />
                            <View style={styles.itemInfoWrapper}>
                                <Text style={styles.itemName}>{item.title}</Text>
                                <Text style={styles.itemDetails}>
                                    {item.quantity} x R${item.price.toFixed(2)}
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
                        onPress={() => setSelectedPayment(method)} // Corrigindo a seleção
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
                    <Text style={styles.totalTxt}>Total: R${getTotal(cartItems)}</Text>
                </View>
                <TouchableOpacity
                    style={styles.checkoutBtn}
                    onPress={handlePlaceOrder}
                >
                    <Text style={styles.checkoutBtnTxt}>Place Order</Text>
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
        fontSize: 18,
        fontWeight: '600',
        color: Colors.black,
        flexWrap: 'wrap',
    },
    checkoutBtn: {
        flex: 1,
        backgroundColor: Colors.primary,
        height: 40,
        justifyContent: 'center',
        alignItems: 'center',
        borderRadius: 5
    },
    checkoutBtnTxt: {
        fontSize: 16,
        fontWeight: '500',
        color: Colors.white
    },
    paymentOption: {
        flexDirection: 'row',
        alignItems: 'center',
        padding: 15,
        marginVertical: 5,
        borderRadius: 8,
        borderWidth: 1,
        borderColor: Colors.lightGray,
        backgroundColor: Colors.white,
    },
    selectedOption: {
        borderColor: Colors.primary,
        backgroundColor: Colors.transparentPrimary,
        borderWidth: 2,
    },
    paymentText: {
        fontSize: 16,
        color: Colors.darkGray,
        marginLeft: 15,
        flex: 1,
    },
    checkmark: {
        marginLeft: 10,
    },
});

export default CheckoutScreen