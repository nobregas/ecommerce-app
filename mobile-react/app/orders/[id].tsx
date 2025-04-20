import { View, Text, ScrollView, StyleSheet, ActivityIndicator, TouchableOpacity } from 'react-native'
import React, { useEffect, useState } from 'react'
import { Stack, useLocalSearchParams, useRouter } from 'expo-router'
import { useHeaderHeight } from '@react-navigation/elements'
import { useOrderStore } from '@/store/orderStore'
import { Colors } from '@/constants/Colors'
import { Image } from 'expo-image'
import { OrderStatus } from '@/service/orderService'
import Ionicons from 'react-native-vector-icons/Ionicons'
import productService from '@/service/productService'

const StatusColors = {
  warning: '#FFA000',
  info: '#2196F3',
  success: '#4CAF50',
  error: '#F44336'
}

interface ProductInfo {
  id: number;
  title: string;
}

const OrderDetails = () => {
  const { id } = useLocalSearchParams()
  const orderId = Number(id)
  const router = useRouter()
  const headerHeight = useHeaderHeight()
  const { fetchOrderWithItems, orderDetails, isLoading, error } = useOrderStore()
  const [productInfos, setProductInfos] = useState<Record<number, ProductInfo>>({})
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (orderId) {
      fetchOrderWithItems(orderId)
    }
  }, [orderId])

  useEffect(() => {
    const fetchProductDetails = async () => {
      if (orderDetails?.items) {
        setLoading(true)
        const productMap: Record<number, ProductInfo> = {}
        
        try {
          const promises = orderDetails.items.map(async (item) => {
            try {
              const product = await productService.getProductDetails(item.productId)
              productMap[item.productId] = {
                id: item.productId,
                title: product.title
              }
            } catch (error) {
              console.error(`Failed to fetch product ${item.productId}:`, error)
              productMap[item.productId] = {
                id: item.productId,
                title: `Product #${item.productId}`
              }
            }
          })
          
          await Promise.all(promises)
          setProductInfos(productMap)
        } catch (error) {
          console.error('Error fetching product details:', error)
        } finally {
          setLoading(false)
        }
      }
    }
    
    if (orderDetails) {
      fetchProductDetails()
    }
  }, [orderDetails])

  const getStatusColor = (status: OrderStatus) => {
    switch (status) {
      case OrderStatus.PENDING:
        return StatusColors.warning
      case OrderStatus.PROCESSING:
        return StatusColors.info
      case OrderStatus.SHIPPED:
        return '#1E88E5'
      case OrderStatus.DELIVERED:
        return StatusColors.success
      case OrderStatus.CANCELLED:
        return StatusColors.error
      default:
        return Colors.gray
    }
  }

  const getStatusLabel = (status: OrderStatus) => {
    switch (status) {
      case OrderStatus.PENDING:
        return 'Pending'
      case OrderStatus.PROCESSING:
        return 'Processing'
      case OrderStatus.SHIPPED:
        return 'Shipped'
      case OrderStatus.DELIVERED:
        return 'Delivered'
      case OrderStatus.CANCELLED:
        return 'Cancelled'
      default:
        return status
    }
  }

  const getPaymentMethodLabel = (method: string) => {
    switch (method) {
      case 'CREDIT_CARD':
        return 'Credit Card'
      case 'DEBIT_CARD':
        return 'Debit Card'
      case 'PIX':
        return 'PIX'
      default:
        return method
    }
  }

  if (isLoading || loading) {
    return (
      <>
        <Stack.Screen
          options={{
            headerShown: true,
            headerTransparent: true,
            headerTitleAlign: 'center',
            title: 'Order Details',
          }}
        />
        <View style={[styles.container, { marginTop: headerHeight }, styles.centerContent]}>
          <ActivityIndicator size="large" color={Colors.primary} />
          <Text style={styles.loadingText}>Loading order details...</Text>
        </View>
      </>
    )
  }

  if (error || !orderDetails) {
    return (
      <>
        <Stack.Screen
          options={{
            headerShown: true,
            headerTransparent: true,
            headerTitleAlign: 'center',
            title: 'Order Details',
          }}
        />
        <View style={[styles.container, { marginTop: headerHeight }, styles.centerContent]}>
          <Ionicons name="alert-circle-outline" size={50} color={StatusColors.error} />
          <Text style={styles.errorText}>Failed to load order details</Text>
          <TouchableOpacity
            style={styles.backButton}
            onPress={() => router.push("/(tabs)")}
          >
            <Text style={styles.backButtonText}>Go to Home</Text>
          </TouchableOpacity>
        </View>
      </>
    )
  }

  return (
    <>
      <Stack.Screen
        options={{
          headerShown: true,
          headerTransparent: true,
          headerTitleAlign: 'center',
          title: `Order #${orderDetails.order.id}`,
        }}
      />
      <ScrollView style={[styles.container, { marginTop: headerHeight }]}>
        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Order Information</Text>
          <View style={styles.infoRow}>
            <Text style={styles.infoLabel}>Number:</Text>
            <Text style={styles.infoValue}>#{orderDetails.order.id}</Text>
          </View>
          <View style={styles.infoRow}>
            <Text style={styles.infoLabel}>Date:</Text>
            <Text style={styles.infoValue}>{orderDetails.order.createdAt}</Text>
          </View>
          <View style={styles.infoRow}>
            <Text style={styles.infoLabel}>Total:</Text>
            <Text style={styles.infoValue}>R${orderDetails.order.totalAmount.toFixed(2)}</Text>
          </View>
          <View style={styles.infoRow}>
            <Text style={styles.infoLabel}>Payment Method:</Text>
            <Text style={styles.infoValue}>{getPaymentMethodLabel(orderDetails.order.paymentMethod)}</Text>
          </View>
          <View style={styles.infoRow}>
            <Text style={styles.infoLabel}>Status:</Text>
            <View style={[styles.statusBadge, { backgroundColor: getStatusColor(orderDetails.order.status) }]}>
              <Text style={styles.statusText}>
                {getStatusLabel(orderDetails.order.status)}
              </Text>
            </View>
          </View>
        </View>

        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Order Items</Text>
          {orderDetails.items.map((item, index) => (
            <View key={index} style={styles.orderItem}>
              <Text style={styles.itemId}>{productInfos[item.productId]?.title || `Product #${item.productId}`}</Text>
              <View style={styles.itemDetails}>
                <Text style={styles.itemQuantity}>{item.quantity}x</Text>
                <Text style={styles.itemPrice}>R${item.price.toFixed(2)}</Text>
              </View>
              <Text style={styles.itemSubtotal}>Subtotal: R${(item.price * item.quantity).toFixed(2)}</Text>
            </View>
          ))}
        </View>
      </ScrollView>
    </>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: Colors.background,
    padding: 20,
  },
  centerContent: {
    justifyContent: 'center',
    alignItems: 'center',
  },
  loadingText: {
    marginTop: 20,
    fontSize: 16,
    color: Colors.gray,
  },
  errorText: {
    marginTop: 10,
    fontSize: 18,
    color: StatusColors.error,
    textAlign: 'center',
  },
  backButton: {
    marginTop: 20,
    backgroundColor: Colors.primary,
    paddingVertical: 10,
    paddingHorizontal: 20,
    borderRadius: 8,
  },
  backButtonText: {
    color: Colors.white,
    fontSize: 16,
    fontWeight: '600',
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
  infoRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  infoLabel: {
    fontSize: 15,
    color: Colors.darkGray,
    fontWeight: '500',
  },
  infoValue: {
    fontSize: 15,
    color: Colors.black,
    fontWeight: '400',
  },
  statusBadge: {
    paddingVertical: 4,
    paddingHorizontal: 10,
    borderRadius: 20,
  },
  statusText: {
    color: Colors.white,
    fontSize: 12,
    fontWeight: '600',
  },
  orderItem: {
    padding: 12,
    backgroundColor: Colors.background,
    borderRadius: 8,
    marginBottom: 10,
  },
  itemId: {
    fontSize: 16,
    fontWeight: '600',
    color: Colors.primary,
    marginBottom: 5,
  },
  itemDetails: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 5,
  },
  itemQuantity: {
    fontSize: 14,
    color: Colors.darkGray,
  },
  itemPrice: {
    fontSize: 14,
    color: Colors.black,
  },
  itemSubtotal: {
    fontSize: 14,
    fontWeight: '500',
    color: Colors.black,
    textAlign: 'right',
  },
})

export default OrderDetails 