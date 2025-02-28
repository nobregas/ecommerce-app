import { FlatList, StyleSheet, Text, Touchable, TouchableOpacity, View } from 'react-native'
import React, { useEffect, useState } from 'react'
import { CartItemType } from '@/types/type'
import { useHeaderHeight } from '@react-navigation/elements'
import axios from 'axios'
import { Stack } from 'expo-router'
import CartItem from '@/components/CartItem'
import { Colors } from '@/constants/Colors'
import Animated, { FadeInDown, SlideInDown } from 'react-native-reanimated'

type Props = {}

const CartScreen = (props: Props) => {

  const headerHeight = useHeaderHeight()
  const [cartItems, setCartItems] = useState<CartItemType[]>([])

  useEffect(() => {
    getCartItems()
  }, [])

  const getCartItems = async () => {
    const URL = `http://10.0.2.2:8000/cart`
    const response = await axios.get(URL)
    setCartItems(response.data)
  }

  return (
    <>
      <Stack.Screen
        options={{ headerShown: true, headerTransparent: true, headerTitleAlign: 'center', title: 'cart' }}
      />
      <View style={[styles.container, { marginTop: headerHeight }]}>
        <FlatList
          data={cartItems}
          keyExtractor={(item) => item.id.toString()}
          renderItem={({ item, index }) => (
            <Animated.View entering={FadeInDown.delay(300 + index * 100).duration(500)}>
              <CartItem item={item} />
            </Animated.View>
          )}
        />
      </View>
      <Animated.View style={styles.footer} entering={SlideInDown.delay(500).duration(500)}>
        <View style={styles.priceInfoWrapper}>
          <Text style={styles.totalTxt}>Total: R$200</Text>
        </View>
        <TouchableOpacity style={styles.checkoutBtn}>
          <Text style={styles.checkoutBtnTxt}>Checkout</Text>
        </TouchableOpacity>
      </Animated.View>
    </>
  )
}

export default CartScreen

const styles = StyleSheet.create({
  container: {
    flex: 1,
    paddingHorizontal: 20
  },
  footer: {
    flexDirection: 'row',
    padding: 20,
    backgroundColor: Colors.white,
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
    height: 40,
    justifyContent: 'center',
    alignItems: 'center',
    borderRadius: 5
  },
  checkoutBtnTxt: {
    fontSize: 16,
    fontWeight: '500',
    color: Colors.white
  }
})