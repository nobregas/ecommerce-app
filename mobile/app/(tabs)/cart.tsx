import { FlatList, StyleSheet, Text, Touchable, TouchableOpacity, View } from 'react-native'
import React, { useEffect, useState } from 'react'
import { CartItemType } from '@/types/type'
import { useHeaderHeight } from '@react-navigation/elements'
import { Link, Stack } from 'expo-router'
import CartItem from '@/components/CartItem'
import { Colors } from '@/constants/Colors'
import Animated, { FadeInDown, SlideInDown } from 'react-native-reanimated'
import { getCartItems, removeFromCart, updateQuantity } from '@/service/ApiService'
import { getTotal } from '@/utils/sharedFunctions'

type Props = {}

const CartScreen = (props: Props) => {

  const headerHeight = useHeaderHeight()
  const [cartItems, setCartItems] = useState<CartItemType[]>([])

  useEffect(() => {
    fetchCartItems()
    const interval = setInterval(fetchCartItems, 5000)

    return () => clearInterval(interval)
  }, [])

  const fetchCartItems = async () => {
    try {
      const cartItemsData = await getCartItems()

      setCartItems(cartItemsData)

    } catch (error) {
      console.log(error)
    }
  }

  const handleUpdateQuantity = async (id: number, newQuantity: number) => {
    const updatedItem = await updateQuantity(id, newQuantity)

    // Update the quantity of an especific cart item
    setCartItems(prevItems => {
      return prevItems.map(item => {
        if (item.id === id) {
          return { ...item, quantity: newQuantity }
        }
        return item
      })
    })
  }

  const handleRemoveItem = async (id: number) => {
    const removedItemId = await removeFromCart(id)

    // Remove the cart item
    setCartItems(prevItems => {
      return prevItems.filter(item => item.id !== id)
    })
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
              <CartItem
                item={item}
                updateQuantity={handleUpdateQuantity}
                removeItem={handleRemoveItem}
              />
            </Animated.View>
          )}
        />
      </View>
      <Animated.View style={styles.footer} entering={SlideInDown.delay(500).duration(500)}>
        <View style={styles.priceInfoWrapper}>
          <Text style={styles.totalTxt}>Total: R${getTotal(cartItems)}</Text>
        </View>
        <Link href="/checkout/index" asChild>
          <TouchableOpacity style={styles.checkoutBtn}>
            <Text style={styles.checkoutBtnTxt}>Checkout</Text>
          </TouchableOpacity>
        </Link>
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